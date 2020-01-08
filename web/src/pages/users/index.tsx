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
import { Box, Button, Flex, Icon } from 'pouncejs';
import { ADMIN_ROLES_ARRAY } from 'Source/constants';
import ListUsersTable from 'Pages/users/subcomponents/list-users-table';
import useSidesheet from 'Hooks/useSidesheet';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import ErrorBoundary from 'Components/error-boundary';

// Parent container for the users management
const UsersContainer: React.FC = () => {
  const { showSidesheet } = useSidesheet();

  return (
    <Box mb={6}>
      <RoleRestrictedAccess allowedRoles={ADMIN_ROLES_ARRAY}>
        <Flex justifyContent="flex-end">
          <Button
            size="large"
            variant="primary"
            onClick={() => showSidesheet({ sidesheet: SIDESHEETS.USER_INVITATION })}
            mb={8}
          >
            <Flex alignItems="center">
              <Icon type="addUser" size="small" mr={2} />
              Invite User
            </Flex>
          </Button>
        </Flex>
      </RoleRestrictedAccess>
      <ErrorBoundary>
        <ListUsersTable />
      </ErrorBoundary>
    </Box>
  );
};

export default UsersContainer;
