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
import { useMutation, gql } from '@apollo/client';
import { Field, Formik } from 'formik';
import { Box, useSnackbar } from 'pouncejs';

import SubmitButton from 'Components/utils/SubmitButton';
import { GET_ORGANIZATION } from 'Pages/general-settings';
import FormikTextInput from 'Components/fields/text-input';
import { extractErrorMessage } from 'Helpers/utils';
import { UpdateOrganizationInput } from 'Generated/schema';

export const UPDATE_ORGANIZATION = gql`
  mutation UpdateCompanyInformation($input: UpdateOrganizationInput!) {
    updateOrganization(input: $input)
  }
`;

interface UpdateCompanyInformationFormValues {
  displayName?: string;
  email?: string;
}

interface ApolloMutationInput {
  input: UpdateOrganizationInput;
}

type UpdateCompanyInformationFormOuterProps = UpdateCompanyInformationFormValues & {
  onSuccess: () => void;
};

export const UpdateCompanyInformationForm: React.FC<UpdateCompanyInformationFormOuterProps> = ({
  displayName,
  email,
  onSuccess,
}) => {
  const { pushSnackbar } = useSnackbar();
  const [
    updateOrganization,
    { loading: updateOrganizationLoading, error: updateOrganizationError, data },
  ] = useMutation<boolean, ApolloMutationInput>(UPDATE_ORGANIZATION);

  React.useEffect(() => {
    if (updateOrganizationError) {
      pushSnackbar({
        variant: 'error',
        title:
          extractErrorMessage(updateOrganizationError) ||
          'Failed to update company information due to an unknown error',
      });
    }
  }, [updateOrganizationError]);

  React.useEffect(() => {
    if (data) {
      pushSnackbar({ variant: 'success', title: `Successfully updated company information` });
      onSuccess();
    }
  }, [data]);

  return (
    <Formik<UpdateCompanyInformationFormValues>
      initialValues={{
        displayName,
        email,
      }}
      onSubmit={async values => {
        await updateOrganization({
          variables: { input: values },
          refetchQueries: [{ query: GET_ORGANIZATION }],
        });
      }}
    >
      {({ handleSubmit }) => (
        <Box>
          <form onSubmit={handleSubmit}>
            <Box mb={8}>
              <Field as={FormikTextInput} name="displayName" label="Name" aria-required />
              <Field as={FormikTextInput} name="email" label="Email" aria-required />
            </Box>
            <SubmitButton
              disabled={updateOrganizationLoading}
              submitting={updateOrganizationLoading}
            >
              Update
            </SubmitButton>
          </form>
        </Box>
      )}
    </Formik>
  );
};

export default UpdateCompanyInformationForm;