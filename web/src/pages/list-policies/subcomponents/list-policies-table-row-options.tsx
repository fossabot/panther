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
import { Dropdown, Icon, IconButton, MenuItem, Box } from 'pouncejs';
import useRouter from 'Hooks/useRouter';
import { PolicySummary } from 'Generated/schema';
import urls from 'Source/urls';
import { READONLY_ROLES_ARRAY } from 'Source/constants';
import useModal from 'Hooks/useModal';
import { MODALS } from 'Components/utils/modal-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';

interface ListPoliciesTableRowOptionsProps {
  policy: PolicySummary;
}

const ListPoliciesTableRowOptions: React.FC<ListPoliciesTableRowOptionsProps> = ({ policy }) => {
  const { history } = useRouter();
  const { showModal } = useModal();

  // @HELP_WANTED
  // The wrapping `<Box>` is needed because of a special reason. You see, the trigger of this
  // Dropdown is added on each row of a table whose rows are clickable. Thus we don't wanna trigger
  // a click on the row each time the Dropdown trigger is clicked. If we had added the
  // `stopPropagation()` on the trigger itself, then the Dropdown wouldn't open.
  // We perhaps can find a better solution to his problem
  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
      <Box onClick={e => e.stopPropagation()}>
        <Dropdown
          trigger={
            <IconButton is="div" variant="default" my={-2}>
              <Icon type="more" size="small" />
            </IconButton>
          }
        >
          <Dropdown.Item onSelect={() => history.push(urls.policies.edit(policy.id))}>
            <MenuItem variant="default">Edit</MenuItem>
          </Dropdown.Item>
          <Dropdown.Item
            onSelect={() =>
              showModal({
                modal: MODALS.DELETE_POLICY,
                props: { policy },
              })
            }
          >
            <MenuItem variant="default">Delete</MenuItem>
          </Dropdown.Item>
        </Dropdown>
      </Box>
    </RoleRestrictedAccess>
  );
};

export default React.memo(ListPoliciesTableRowOptions);
