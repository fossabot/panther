import React from 'react';
import { Button, ButtonProps, useSnackbar } from 'pouncejs';
import { useMutation, gql } from '@apollo/client';

import { READONLY_ROLES_ARRAY } from 'Source/constants';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import { RESOURCE_DETAILS } from 'Pages/resource-details';
import { POLICY_DETAILS } from 'Pages/policy-details';
import { ResourceDetails, RemediateResourceInput, PolicyDetails } from 'Generated/schema';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import { extractErrorMessage } from 'Helpers/utils';

interface RemediationButtonProps {
  buttonVariant: ButtonProps['variant'];
  resourceId: ResourceDetails['id'];
  policyId: PolicyDetails['id'];
}

const REMEDIATE_RESOURCE = gql`
  mutation RemediateResource($input: RemediateResourceInput!) {
    remediateResource(input: $input)
  }
`;

interface ApolloMutationInput {
  input: RemediateResourceInput;
}

const RemediationButton: React.FC<RemediationButtonProps> = ({
  buttonVariant,
  resourceId,
  policyId,
}) => {
  const { pushSnackbar } = useSnackbar();

  // Prepare the remediation mutation.
  const [
    remediateResource,
    { data: remediationSuccess, error: remediationError, loading: remediationInProgress },
  ] = useMutation<boolean, ApolloMutationInput>(REMEDIATE_RESOURCE, {
    mutation: REMEDIATE_RESOURCE,
    awaitRefetchQueries: true,
    refetchQueries: [getOperationName(RESOURCE_DETAILS), getOperationName(POLICY_DETAILS)],
    variables: {
      input: {
        resourceId,
        policyId,
      },
    },
  });

  React.useEffect(() => {
    if (remediationError) {
      pushSnackbar({
        variant: 'error',
        title: extractErrorMessage(remediationError) || 'Failed to apply remediation',
      });
    }
  }, [remediationError]);

  React.useEffect(() => {
    if (remediationSuccess) {
      pushSnackbar({ variant: 'success', title: 'Remediation has been applied successfully' });
    }
  }, [remediationSuccess]);

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
      <Button
        size="small"
        variant={buttonVariant}
        onClick={e => {
          // Table row is clickable, we don't want to navigate away
          e.stopPropagation();
          remediateResource();
        }}
        disabled={remediationInProgress}
      >
        Remediate
      </Button>
    </RoleRestrictedAccess>
  );
};

export default React.memo(RemediationButton);
