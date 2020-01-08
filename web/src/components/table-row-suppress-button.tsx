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
