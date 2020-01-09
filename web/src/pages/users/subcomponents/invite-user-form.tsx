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
import * as Yup from 'yup';

import { useMutation, gql } from '@apollo/client';
import { Field, Formik } from 'formik';
import { Alert, Box, Flex, useSnackbar } from 'pouncejs';
import { RoleNameEnum, InviteUserInput } from 'Generated/schema';
import { LIST_USERS } from 'Pages/users/subcomponents/list-users-table';
import SubmitButton from 'Components/utils/SubmitButton';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import FormikTextInput from 'Components/fields/text-input';
import FormikCombobox from 'Components/fields/combobox';
import { extractErrorMessage } from 'Helpers/utils';

const INVITE_USER = gql`
  mutation InviteUser($input: InviteUserInput!) {
    inviteUser(input: $input) {
      id
    }
  }
`;

interface ApolloMutationInput {
  input: InviteUserInput;
}

interface InviteUserFormValues {
  email?: string;
  familyName?: string;
  givenName?: string;
  role?: RoleNameEnum;
}

interface InviteUserFormProps {
  onSuccess: () => void;
}

const initialValues = {
  email: '',
  familyName: '',
  givenName: '',
  role: RoleNameEnum.Analyst,
};

const validationSchema = Yup.object().shape({
  email: Yup.string().required('Email is required'),
  familyName: Yup.string().required('Last name is required'),
  givenName: Yup.string().required('First name is required'),
  role: Yup.string().required('Role is required'),
});

const roleValues = Object.values(RoleNameEnum);

export const InviteUserForm: React.FC<InviteUserFormProps> = ({ onSuccess }) => {
  const [inviteUser, { error: inviteUserError, data }] = useMutation<boolean, ApolloMutationInput>(
    INVITE_USER
  );
  const { pushSnackbar } = useSnackbar();

  React.useEffect(() => {
    if (data) {
      pushSnackbar({ variant: 'success', title: `Successfully invited user` });
      onSuccess();
    }
  }, [data]);

  return (
    <Formik<InviteUserFormValues>
      validationSchema={validationSchema}
      initialValues={initialValues}
      onSubmit={async values => {
        await inviteUser({
          variables: {
            input: {
              email: values.email,
              familyName: values.familyName,
              givenName: values.givenName,
              role: values.role,
            },
          },
          // TODO: Find a better way to update the cache using response from invite user
          refetchQueries: [getOperationName(LIST_USERS)],
        });
      }}
    >
      {({ handleSubmit, isSubmitting, dirty, isValid }) => (
        <form onSubmit={handleSubmit}>
          {inviteUserError && (
            <Alert
              variant="error"
              title="Failed to invite user"
              description={
                extractErrorMessage(inviteUserError) ||
                'Failed to invite user due to an unforeseen error'
              }
              mb={6}
            />
          )}
          <Box mb={8}>
            <Flex justifyContent="space-between">
              <Field name="givenName" as={FormikTextInput} label="First Name" />
              <Field name="familyName" as={FormikTextInput} label="Family Name" />
            </Flex>
            <Field name="email" as={FormikTextInput} type="email" label="Email" />
            <Field name="role" as={FormikCombobox} label="Role" items={roleValues} />
          </Box>
          <SubmitButton
            width={1}
            disabled={isSubmitting || !isValid || !dirty}
            submitting={isSubmitting}
          >
            Invite User
          </SubmitButton>
        </form>
      )}
    </Formik>
  );
};

export default InviteUserForm;