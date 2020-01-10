package mage

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
)

const (
	keysDirectory        = "keys"
	certificateFile      = keysDirectory + "/panther-tls-public.crt"
	privateKeyFile       = keysDirectory + "/panther-tls-private.key"
	keyLength            = 2048
	certFilePermissions  = 0700
	certificateOutputKey = "WebApplicationCertificateArn"
)

// Upload a local self-signed TLS certificate to ACM. Only needs to happen once per installation
func uploadLocalCertificate(awsSession *session.Session) (string, error) {
	// Check if certificate has already been uploaded
	certArn, err := getExistingCertificate(awsSession)
	if err != nil {
		return "", err
	}
	if certArn != "" {
		fmt.Println("deploy: ACM certificate already exists")
		return certArn, nil
	}
	fmt.Println("deploy: uploading ACM certificate")

	// Ensure the certificate and key file exist. If not, create them.
	_, certErr := os.Stat(certificateFile)
	_, keyErr := os.Stat(certificateFile)
	if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
		if err := generateKeys(); err != nil {
			return "", err
		}
	}

	certificateFile, certificateFileErr := os.Open(certificateFile)
	if certificateFileErr != nil {
		return "", certificateFileErr
	}
	defer func() { _ = certificateFile.Close() }()

	privateKeyFile, privateKeyFileErr := os.Open(privateKeyFile)
	if privateKeyFileErr != nil {
		return "", privateKeyFileErr
	}
	defer func() { _ = privateKeyFile.Close() }()

	certificateBytes, err := ioutil.ReadAll(certificateFile)
	if err != nil {
		return "", err
	}
	privateKeyBytes, err := ioutil.ReadAll(privateKeyFile)
	if err != nil {
		return "", err
	}

	input := &acm.ImportCertificateInput{
		Certificate: certificateBytes,
		PrivateKey:  privateKeyBytes,
		Tags: []*acm.Tag{
			{
				Key:   aws.String("Application"),
				Value: aws.String("Panther"),
			},
		},
	}

	acmClient := acm.New(awsSession)
	output, err := acmClient.ImportCertificate(input)
	if err != nil {
		return "", err
	}
	return *output.CertificateArn, nil
}

func getExistingCertificate(awsSession *session.Session) (string, error) {
	outputs, err := getStackOutputs(awsSession, applicationStack)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != "ValidationError" || !strings.HasSuffix(awsErr.Code(), "does not exist") {
				return "", nil
			}
		}
		return "", err
	}
	if arn, ok := outputs[certificateOutputKey]; ok {
		return arn, nil
	}
	return "", nil
}

// Generate the self signed private key and certificate for HTTPS access to the web application
func generateKeys() error {
	fmt.Println("deploy: WARNING no ACM certificate ARN provided and no certificate file provided, generating a self-signed certificate")
	// Create the private key
	key, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}

	// Setup the certificate template
	certificateTemplate := x509.Certificate{
		BasicConstraintsValid: true,
		// AWS will not attach a certificate that does not have a domain specified
		// example.com is reserved by IANA and is not available for registration so there is no risk
		// of confusion about us trying to MITM someone (ref: https://www.iana.org/domains/reserved)
		DNSNames:     []string{"example.com"},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:     time.Now().Add(time.Hour * 24 * 365),
		NotBefore:    time.Now(),
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Panther User"},
		},
	}

	// Create the certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &key.PublicKey, key)
	if err != nil {
		return err
	}

	// Create the keys directory if it does not already exist
	err = os.MkdirAll(keysDirectory, certFilePermissions)
	if err != nil {
		return err
	}

	// PEM encode the certificate and write it to disk
	certBuffer := &bytes.Buffer{}
	err = pem.Encode(
		certBuffer,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes},
	)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(certificateFile, certBuffer.Bytes(), certFilePermissions)
	if err != nil {
		return err
	}

	// PEM Encode the private key and write it to disk
	keyBuffer := &bytes.Buffer{}
	err = pem.Encode(
		keyBuffer,
		&pem.Block{Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(privateKeyFile, keyBuffer.Bytes(), certFilePermissions)
}
