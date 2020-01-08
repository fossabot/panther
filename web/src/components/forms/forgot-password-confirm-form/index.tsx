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

import React from 'react';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import { Alert, Box, useSnackbar } from 'pouncejs';
import SubmitButton from 'Components/utils/SubmitButton';
import FormikTextInput from 'Components/fields/text-input';
import useRouter from 'Hooks/useRouter';
import useAuth from 'Hooks/useAuth';
import urls from 'Source/urls';

interface ForgotPasswordConfirmFormProps {
  email: string;
  token: string;
}

interface ForgotPasswordConfirmFormValues {
  newPassword: string;
  confirmNewPassword: string;
}

const validationSchema = Yup.object().shape({
  newPassword: Yup.string().required(),
  confirmNewPassword: Yup.string()
    .oneOf([Yup.ref('newPassword')], 'Passwords must match')
    .required(),
});

const ForgotPasswordConfirmForm: React.FC<ForgotPasswordConfirmFormProps> = ({ email, token }) => {
  const { history } = useRouter();
  const { resetPassword } = useAuth();
  const { pushSnackbar } = useSnackbar();

  const initialValues = {
    newPassword: '',
    confirmNewPassword: '',
  };

  return (
    <Formik<ForgotPasswordConfirmFormValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={({ newPassword }, { setStatus }) =>
        resetPassword({
          email,
          token,
          newPassword,
          onSuccess: () => {
            pushSnackbar({ variant: 'info', title: 'Password changed successfully!' });
            history.replace(urls.account.auth.signIn());
          },
          onError: ({ message }) =>
            setStatus({
              title: 'Houston, we have a problem',
              message,
            }),
        })
      }
    >
      {({ isValid, handleSubmit, isSubmitting, status, dirty }) => (
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
            autoco
            aria-required
            autoComplete="new-password"
            mb={6}
          />
          <Field
            as={FormikTextInput}
            label="Confirm New Password"
            placeholder="Type your new password again..."
            type="password"
            name="confirmNewPassword"
            aria-required
            autoComplete="new-password"
            mb={6}
          />
          <SubmitButton
            width={1}
            submitting={isSubmitting}
            disabled={isSubmitting || !isValid || !dirty}
          >
            Update password
          </SubmitButton>
        </Box>
      )}
    </Formik>
  );
};

export default ForgotPasswordConfirmForm;
