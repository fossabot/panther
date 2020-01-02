import React from 'react';
import Panel from 'Components/panel';
import { Alert, Button, Card, Box } from 'pouncejs';
import PolicyForm from 'Components/forms/policy-form';
import { GetPolicyInput, PolicyDetails } from 'Generated/schema';
import useModal from 'Hooks/useModal';
import { useMutation, useQuery, gql } from '@apollo/client';
import useRouter from 'Hooks/useRouter';
import TablePlaceholder from 'Components/table-placeholder';
import { MODALS } from 'Components/utils/modal-context';
import { READONLY_ROLES_ARRAY } from 'Source/constants';
import useEditRule from 'Hooks/useEditRule';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import Page403 from 'Pages/403';
import { extractErrorMessage } from 'Helpers/utils';

const POLICY_DETAILS = gql`
  query PolicyDetails($input: GetPolicyInput!) {
    policy(input: $input) {
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

const UPDATE_POLICY = gql`
  mutation UpdatePolicy($input: CreateOrModifyPolicyInput!) {
    updatePolicy(input: $input) {
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

interface ApolloQueryData {
  policy: PolicyDetails;
}

interface ApolloQueryInput {
  input: GetPolicyInput;
}

interface ApolloMutationData {
  updatePolicy: PolicyDetails;
}

interface ApolloMutationInput {
  input: GetPolicyInput;
}

const EditPolicyPage: React.FC = () => {
  const { match } = useRouter<{ id: string }>();
  const { showModal } = useModal();

  const { error: fetchPolicyError, data: queryData, loading: isFetchingPolicy } = useQuery<
    ApolloQueryData,
    ApolloQueryInput
  >(POLICY_DETAILS, {
    fetchPolicy: 'cache-and-network',
    variables: {
      input: {
        policyId: match.params.id,
      },
    },
  });

  const mutation = useMutation<ApolloMutationData, ApolloMutationInput>(UPDATE_POLICY);

  const { initialValues, handleSubmit, error: updateError } = useEditRule<ApolloMutationData>({
    mutation,
    type: 'policy',
    rule: queryData?.policy,
  });

  if (isFetchingPolicy) {
    return (
      <Card p={9}>
        <TablePlaceholder rowCount={5} rowHeight={15} />
        <TablePlaceholder rowCount={1} rowHeight={100} />
      </Card>
    );
  }

  if (fetchPolicyError) {
    return (
      <Alert
        mb={6}
        variant="error"
        title="Couldn't load the policy details"
        description={
          extractErrorMessage(fetchPolicyError) ||
          'There was an error when performing your request, please contact support@runpanther.io'
        }
      />
    );
  }

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY} fallback={<Page403 />}>
      <Box mb={6}>
        <Panel
          size="large"
          title="Policy Settings"
          actions={
            <Button
              variant="default"
              size="large"
              color="red300"
              onClick={() =>
                showModal({
                  modal: MODALS.DELETE_POLICY,
                  props: { policy: queryData.policy },
                })
              }
            >
              Delete
            </Button>
          }
        >
          <PolicyForm initialValues={initialValues} onSubmit={handleSubmit} />
        </Panel>
        {updateError && (
          <Alert
            mt={2}
            mb={6}
            variant="error"
            title={
              extractErrorMessage(updateError) ||
              'Unknown error occured during update. Please contact support@runpanther.io'
            }
          />
        )}
      </Box>
    </RoleRestrictedAccess>
  );
};

export default EditPolicyPage;
