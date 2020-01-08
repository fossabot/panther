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

import { Alert, Box, Text, Flex } from 'pouncejs';
import { Field, Formik } from 'formik';
import QRCode from 'qrcode.react';
import * as React from 'react';
import * as Yup from 'yup';
import { formatSecretCode } from 'Helpers/utils';
import SubmitButton from 'Components/utils/SubmitButton';
import FormikTextInput from 'Components/fields/text-input';
import useAuth from 'Hooks/useAuth';

interface TotpFormValues {
  mfaCode: string;
}

const initialValues = {
  mfaCode: '',
};

const validationSchema = Yup.object().shape({
  mfaCode: Yup.string()
    .matches(/\b\d{6}\b/, 'Code should contain exactly six digits.')
    .required(),
});

export const TotpForm: React.FC = () => {
  const [code, setCode] = React.useState('');
  const { userInfo, verifyTotpSetup, requestTotpSecretCode } = useAuth();

  React.useEffect(() => {
    (async () => {
      setCode(await requestTotpSecretCode());
    })();
  }, []);

  return (
    <Formik<TotpFormValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={async ({ mfaCode }, { setStatus }) =>
        verifyTotpSetup({
          mfaCode,
          onError: ({ message }) =>
            setStatus({
              title: 'Authentication failed',
              message,
            }),
        })
      }
    >
      {({ handleSubmit, isSubmitting, status, isValid, dirty }) => (
        <Box is="form" width="100%" onSubmit={handleSubmit}>
          {status && (
            <Alert variant="error" title={status.title} description={status.message} mb={6} />
          )}
          <Flex justifyContent="center" mb={6} width={1}>
            <QRCode value={formatSecretCode(code, userInfo.email)} />
          </Flex>
          <Field
            autoFocus
            as={FormikTextInput}
            placeholder="The 6-digit MFA code"
            name="mfaCode"
            autoComplete="off"
            aria-required
            mb={6}
          />
          <SubmitButton
            width={1}
            submitting={isSubmitting}
            disabled={isSubmitting || !isValid || !dirty}
          >
            Verify
          </SubmitButton>
          <Text color="grey200" size="small" mt={10} textAlign="center">
            Open any two-factor authentication app, scan the barcode and then enter the MFA code to
            complete the sign-in. Popular software options include{' '}
            <a
              href="https://duo.com/product/trusted-users/two-factor-authentication/duo-mobile"
              target="_blank"
              rel="noopener noreferrer"
            >
              Duo
            </a>
            ,{' '}
            <a
              href="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en"
              target="_blank"
              rel="noopener noreferrer"
            >
              Google authenticator
            </a>
            ,{' '}
            <a
              href="https://lastpass.com/misc_download2.php"
              target="_blank"
              rel="noopener noreferrer"
            >
              LastPass
            </a>{' '}
            and{' '}
            <a
              href="https://1password.com/downloads/mac/"
              target="_blank"
              rel="noopener noreferrer"
            >
              1Password
            </a>
            .
          </Text>
        </Box>
      )}
    </Formik>
  );
};

export default TotpForm;
