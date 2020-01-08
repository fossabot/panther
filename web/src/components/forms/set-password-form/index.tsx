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

import { Field, Formik } from 'formik';
import React from 'react';
import * as Yup from 'yup';
import { createYupPasswordValidationSchema } from 'Helpers/utils';
import { Alert, Box, Text } from 'pouncejs';
import SubmitButton from 'Components/utils/SubmitButton';
import FormikTextInput from 'Components/fields/text-input';
import useAuth from 'Hooks/useAuth';

interface SetPasswordFormValues {
  confirmNewPassword: string;
  newPassword: string;
  formErrors?: string[];
}

const initialValues = {
  confirmNewPassword: '',
  newPassword: '',
};

const validationSchema = Yup.object().shape({
  confirmNewPassword: createYupPasswordValidationSchema().oneOf(
    [Yup.ref('newPassword')],
    'Passwords must match'
  ),
  newPassword: createYupPasswordValidationSchema(),
});

const SetPasswordForm: React.FC = () => {
  const { setNewPassword } = useAuth();

  return (
    <Formik<SetPasswordFormValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={async ({ newPassword }, { setStatus }) =>
        setNewPassword({
          newPassword,
          onError: ({ message }) =>
            setStatus({
              title: 'Update password failed',
              message,
            }),
        })
      }
    >
      {({ handleSubmit, status, isSubmitting, isValid, dirty }) => (
        <Box is="form" width={1} onSubmit={handleSubmit}>
          {status && (
            <Alert variant="error" title={status.title} description={status.message} mb={6} />
          )}
          <Field
            as={FormikTextInput}
            label="New Password"
            placeholder="Type your new password..."
            type="password"
            name="newPassword"
            aria-required
            mb={6}
          />
          <Field
            as={FormikTextInput}
            label="Confirm New Password"
            placeholder="Your new password again..."
            type="password"
            name="confirmNewPassword"
            aria-required
            mb={6}
          />
          <SubmitButton
            width={1}
            submitting={isSubmitting}
            disabled={isSubmitting || !isValid || !dirty}
          >
            Set password
          </SubmitButton>
          <Text size="medium" mt={6} color="grey200">
            By continuing, you agree to Panther&apos;s&nbsp;
            <a
              href="https://panther-public-shared-assets.s3-us-west-2.amazonaws.com/EULA.pdf"
              target="_blank"
              rel="noopener noreferrer"
            >
              End User License Agreement
            </a>{' '}
            and acknowledge you have read the&nbsp;
            <a
              href="https://panther-public-shared-assets.s3-us-west-2.amazonaws.com/PrivacyPolicy.pdf"
              target="_blank"
              rel="noopener noreferrer"
            >
              Privacy Policy
            </a>
          </Text>
        </Box>
      )}
    </Formik>
  );
};

export default SetPasswordForm;
