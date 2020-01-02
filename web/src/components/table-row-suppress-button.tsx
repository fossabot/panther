import React from 'react';
import { Button, ButtonProps } from 'pouncejs';
import { ResourceDetails, PolicyDetails } from 'Generated/schema';
import usePolicySuppression from 'Hooks/usePolicySuppression';
import { READONLY_ROLES_ARRAY } from 'Source/constants';
import RoleRestrictedAccess from 'Components/role-restricted-access';

interface SuppressButtonProps {
  buttonVariant: ButtonProps['variant'];
  resourcePatterns: ResourceDetails['id'][];
  policyIds: PolicyDetails['id'][];
}

const SuppressButton: React.FC<SuppressButtonProps> = ({
  buttonVariant,
  policyIds,
  resourcePatterns,
}) => {
  const { suppressPolicies, loading } = usePolicySuppression({ policyIds, resourcePatterns });

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
      <Button
        size="small"
        variant={buttonVariant}
        onClick={e => {
          // Table row is clickable, we don't want to navigate away
          e.stopPropagation();
          suppressPolicies();
        }}
        disabled={loading}
      >
        Ignore
      </Button>
    </RoleRestrictedAccess>
  );
};

export default React.memo(SuppressButton);
