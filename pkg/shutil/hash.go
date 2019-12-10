package shutil

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// SHA256 computes the hash of the given file.
func SHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
