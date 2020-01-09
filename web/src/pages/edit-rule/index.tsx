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
import { Alert, Button, Card, Box } from 'pouncejs';
import RuleForm from 'Components/forms/rule-form';
import { GetRuleInput, RuleDetails } from 'Generated/schema';
import useModal from 'Hooks/useModal';
import { READONLY_ROLES_ARRAY } from 'Source/constants';
import { useMutation, useQuery, gql } from '@apollo/client';
import useRouter from 'Hooks/useRouter';
import TablePlaceholder from 'Components/table-placeholder';
import { MODALS } from 'Components/utils/modal-context';
import useEditRule from 'Hooks/useEditRule';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import Page403 from 'Pages/403';
import { extractErrorMessage } from 'Helpers/utils';

const RULE_DETAILS = gql`
  query RuleDetails($input: GetRuleInput!) {
    rule(input: $input) {
      description
      displayName
      enabled
      id
      reference
      logTypes
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

const UPDATE_RULE = gql`
  mutation UpdateRule($input: CreateOrModifyRuleInput!) {
    updateRule(input: $input) {
      description
      displayName
      enabled
      id
      reference
      logTypes
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

interface ApolloQueryData {
  rule: RuleDetails;
}

interface ApolloQueryInput {
  input: GetRuleInput;
}

interface ApolloMutationData {
  updateRule: RuleDetails;
}

interface ApolloMutationInput {
  input: GetRuleInput;
}

const EditRulePage: React.FC = () => {
  const { match } = useRouter<{ id: string }>();
  const { showModal } = useModal();

  const { error: fetchRuleError, data: queryData, loading: isFetchingRule } = useQuery<
    ApolloQueryData,
    ApolloQueryInput
  >(RULE_DETAILS, {
    fetchPolicy: 'cache-and-network',
    variables: {
      input: {
        ruleId: match.params.id,
      },
    },
  });

  const mutation = useMutation<ApolloMutationData, ApolloMutationInput>(UPDATE_RULE);

  const { initialValues, handleSubmit, error: updateError } = useEditRule<ApolloMutationData>({
    mutation,
    type: 'rule',
    rule: queryData?.rule,
  });

  if (isFetchingRule) {
    return (
      <Card p={9}>
        <TablePlaceholder rowCount={5} rowHeight={15} />
        <TablePlaceholder rowCount={1} rowHeight={100} />
      </Card>
    );
  }

  if (fetchRuleError) {
    return (
      <Alert
        mb={6}
        variant="error"
        title="Couldn't load the rule details"
        description={
          extractErrorMessage(fetchRuleError) ||
          'There was an error when performing your request, please contact support@runpanther.io'
        }
      />
    );
  }

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY} fallback={<Page403 />}>
      <Box mb={10}>
        <Panel
          size="large"
          title="Rule Settings"
          actions={
            <Button
              variant="default"
              size="large"
              color="red300"
              onClick={() =>
                showModal({
                  modal: MODALS.DELETE_RULE,
                  props: { rule: queryData.rule },
                })
              }
            >
              Delete
            </Button>
          }
        >
          <RuleForm initialValues={initialValues} onSubmit={handleSubmit} />
        </Panel>
        {updateError && (
          <Alert
            mt={2}
            mb={6}
            variant="error"
            title={
              extractErrorMessage(updateError) ||
              'An unknown error occured as were trying to update your rule'
            }
          />
        )}
      </Box>
    </RoleRestrictedAccess>
  );
};

export default EditRulePage;