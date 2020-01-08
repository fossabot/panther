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
import Panel from 'Components/panel';
import { Alert, Box } from 'pouncejs';
import urls from 'Source/urls';
import PolicyForm from 'Components/forms/policy-form';
import { GetPolicyInput, PolicyDetails, ResourceDetails } from 'Generated/schema';
import { useMutation, gql } from '@apollo/client';
import { DEFAULT_POLICY_FUNCTION, READONLY_ROLES_ARRAY } from 'Source/constants';
import useCreateRule from 'Hooks/useCreateRule';
import { LIST_POLICIES } from 'Pages/list-policies';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import Page403 from 'Pages/403';
import { extractErrorMessage } from 'Helpers/utils';

const initialValues: PolicyDetails = {
  autoRemediationId: '',
  autoRemediationParameters: '{}',
  description: '',
  displayName: '',
  enabled: true,
  suppressions: [],
  id: '',
  reference: '',
  resourceTypes: [],
  runbook: '',
  severity: null,
  tags: [],
  body: DEFAULT_POLICY_FUNCTION,
  tests: [],
};

const CREATE_POLICY = gql`
  mutation CreatePolicy($input: CreateOrModifyPolicyInput!) {
    addPolicy(input: $input) {
      autoRemediationId
      autoRemediationParameters
      description
      displayName
      enabled
      suppressions
      id
      reference
      resourceTypes
      runbook
      severity
      tags
      body
      tests {
        expectedResult
        name
        resource
        resourceType
      }
    }
  }
`;

interface ApolloMutationData {
  addPolicy: ResourceDetails;
}

interface ApolloMutationInput {
  input: GetPolicyInput;
}

const EditPolicyPage: React.FC = () => {
  const mutation = useMutation<ApolloMutationData, ApolloMutationInput>(CREATE_POLICY, {
    refetchQueries: [getOperationName(LIST_POLICIES)],
  });

  const { handleSubmit, error } = useCreateRule<ApolloMutationData>({
    mutation,
    getRedirectUri: data => urls.policies.details(data.addPolicy.id),
  });

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY} fallback={<Page403 />}>
      <Box mb={6}>
        <Panel size="large" title="Policy Settings">
          <PolicyForm initialValues={initialValues} onSubmit={handleSubmit} />
        </Panel>
        {error && (
          <Alert
            mt={2}
            mb={6}
            variant="error"
            title={
              extractErrorMessage(error) ||
              'An unknown error occured as we were trying to create your policy'
            }
          />
        )}
      </Box>
    </RoleRestrictedAccess>
  );
};

export default EditPolicyPage;
