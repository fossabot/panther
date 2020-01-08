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
import urls from 'Source/urls';
import { Button, Dropdown, Flex, Icon, MenuItem } from 'pouncejs';
import useRouter from 'Hooks/useRouter';
import useSidesheet from 'Hooks/useSidesheet';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import { READONLY_ROLES_ARRAY } from 'Source/constants';

const CreateButton: React.FC = () => {
  const { history } = useRouter();
  const { showSidesheet } = useSidesheet();

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
      <Dropdown
        width={1}
        trigger={
          <Button size="large" variant="primary" is="div">
            <Flex>
              <Icon type="add" size="small" mr={2} />
              Create new
            </Flex>
          </Button>
        }
      >
        <Dropdown.Item onSelect={() => history.push(urls.rules.create())}>
          <MenuItem variant="default">Single</MenuItem>
        </Dropdown.Item>
        <Dropdown.Item
          onSelect={() =>
            showSidesheet({
              sidesheet: SIDESHEETS.POLICY_BULK_UPLOAD,
              props: { type: 'rule' },
            })
          }
        >
          <MenuItem variant="default">Bulk Upload</MenuItem>
        </Dropdown.Item>
      </Dropdown>
    </RoleRestrictedAccess>
  );
};

export default CreateButton;
