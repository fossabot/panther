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

import * as React from 'react';
import { Alert, Box, Flex, useSnackbar } from 'pouncejs';
import { Field, Formik } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import SubmitButton from 'Components/utils/SubmitButton';
import useAuth from 'Hooks/useAuth';

interface EditProfileFormProps {
  onSuccess: () => void;
}

interface EditProfileFormValues {
  givenName: string;
  familyName: string;
  email: string;
}

const EditProfileForm: React.FC<EditProfileFormProps> = ({ onSuccess }) => {
  const { userInfo, updateUserInfo } = useAuth();
  const { pushSnackbar } = useSnackbar();

  const initialValues = {
    email: userInfo.email || '',
    familyName: userInfo.family_name || '',
    givenName: userInfo.given_name || '',
  };

  return (
    <Formik<EditProfileFormValues>
      initialValues={initialValues}
      onSubmit={async ({ givenName, familyName }, { setStatus }) =>
        updateUserInfo({
          newAttributes: {
            given_name: givenName,
            family_name: familyName,
          },
          onSuccess: () => {
            onSuccess();
            pushSnackbar({ title: 'Successfully updated profile!', variant: 'success' });
          },
          onError: ({ message }) =>
            setStatus({
              title: 'Unable to update profile',
              message,
            }),
        })
      }
    >
      {({ handleSubmit, status, isSubmitting, isValid, dirty }) => (
        <Box is="form" onSubmit={handleSubmit}>
          {status && (
            <Alert variant="error" title={status.title} description={status.message} mb={6} />
          )}
          <Field
            as={FormikTextInput}
            label="Email address"
            placeholder="john@doe.com"
            disabled
            name="email"
            aria-required
            readonly
            mb={3}
          />
          <Flex mb={6} justifyContent="space-between">
            <Field
              as={FormikTextInput}
              label="First Name"
              placeholder="John"
              name="givenName"
              aria-required
            />
            <Field
              as={FormikTextInput}
              label="Last Name"
              placeholder="Doe"
              name="familyName"
              aria-required
            />
          </Flex>
          <SubmitButton
            width={1}
            submitting={isSubmitting}
            disabled={isSubmitting || !isValid || !dirty}
          >
            Update
          </SubmitButton>
        </Box>
      )}
    </Formik>
  );
};

export default EditProfileForm;
