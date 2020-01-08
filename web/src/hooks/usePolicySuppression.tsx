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
import { PolicyDetails, SuppressPoliciesInput, ResourceDetails } from 'Generated/schema';
import { useMutation, gql } from '@apollo/client';
import { useSnackbar } from 'pouncejs';
import { RESOURCE_DETAILS } from 'Pages/resource-details';
import { POLICY_DETAILS } from 'Pages/policy-details';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import { extractErrorMessage } from 'Helpers/utils';

const SUPPRESS_POLICIES = gql`
  mutation SuppressPolicy($input: SuppressPoliciesInput!) {
    suppressPolicies(input: $input)
  }
`;

interface ApolloMutationInput {
  input: SuppressPoliciesInput;
}

interface UsePolicySuppressionProps {
  /** A list of IDs whose corresponding policies should receive the suppression */
  policyIds: PolicyDetails['id'][];

  /** A list of resource patterns (globs) whose matching resources should neglect the above policies
   * during their checks. In other words the resource patterns that should be suppressed for the
   * above policies
   */
  resourcePatterns: ResourceDetails['id'][];
}
const usePolicySuppression = ({ policyIds, resourcePatterns }: UsePolicySuppressionProps) => {
  const [suppressPolicies, { data, loading, error }] = useMutation<boolean, ApolloMutationInput>(
    SUPPRESS_POLICIES,
    {
      awaitRefetchQueries: true,
      refetchQueries: [getOperationName(RESOURCE_DETAILS), getOperationName(POLICY_DETAILS)],
      variables: {
        input: { policyIds, resourcePatterns },
      },
    }
  );

  const { pushSnackbar } = useSnackbar();
  React.useEffect(() => {
    if (error) {
      pushSnackbar({
        variant: 'error',
        title:
          extractErrorMessage(error) ||
          'Failed to apply suppression due to an unknown and unpredicted error',
      });
    }
  }, [error]);

  React.useEffect(() => {
    if (data) {
      pushSnackbar({ variant: 'success', title: 'Suppression applied successfully' });
    }
  }, [data]);

  return React.useMemo(() => ({ suppressPolicies, data, loading, error }), [
    suppressPolicies,
    data,
    loading,
    error,
  ]);
};

export default usePolicySuppression;
