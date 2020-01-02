import React from 'react';
import Panel from 'Components/panel';
import { Alert, Box } from 'pouncejs';
import urls from 'Source/urls';
import RuleForm from 'Components/forms/rule-form';
import { GetRuleInput, RuleDetails } from 'Generated/schema';

import { useMutation, gql } from '@apollo/client';
import { DEFAULT_RULE_FUNCTION, READONLY_ROLES_ARRAY } from 'Source/constants';
import useCreateRule from 'Hooks/useCreateRule';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import Page403 from 'Pages/403';
import { extractErrorMessage } from 'Helpers/utils';

const initialValues: RuleDetails = {
  description: '',
  displayName: '',
  enabled: true,
  id: '',
  reference: '',
  logTypes: [],
  runbook: '',
  severity: null,
  tags: [],
  body: DEFAULT_RULE_FUNCTION,
  tests: [],
};

const CREATE_RULE = gql`
  mutation CreateRule($input: CreateOrModifyRuleInput!) {
    addRule(input: $input) {
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

interface ApolloMutationData {
  addRule: RuleDetails;
}

interface ApolloMutationInput {
  input: GetRuleInput;
}

const CreateRulePage: React.FC = () => {
  const mutation = useMutation<ApolloMutationData, ApolloMutationInput>(CREATE_RULE);

  const { handleSubmit, error } = useCreateRule<ApolloMutationData>({
    mutation,
    getRedirectUri: data => urls.rules.details(data.addRule.id),
  });

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY} fallback={<Page403 />}>
      <Box mb={10}>
        <Panel size="large" title="Rule Settings">
          <RuleForm initialValues={initialValues} onSubmit={handleSubmit} />
        </Panel>
        {error && (
          <Alert
            mt={2}
            mb={6}
            variant="error"
            title={
              extractErrorMessage(error) ||
              'An unknown error occured as we were trying to create your rule'
            }
          />
        )}
      </Box>
    </RoleRestrictedAccess>
  );
};

export default CreateRulePage;
