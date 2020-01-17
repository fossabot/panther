/**
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program [The enterprise software] is licensed under the terms of a commercial license
 * available from Panther Labs Inc ("Panther Commercial License") by contacting contact@runpanther.com.
 * All use, distribution, and/or modification of this software, whether commercial or non-commercial,
 * falls under the Panther Commercial License to the extent it is permitted.
 */

import React from 'react';
import { Box, Button, Flex, Icon } from 'pouncejs';
import { ADMIN_ROLES_ARRAY } from 'Source/constants';
import useSidesheet from 'Hooks/useSidesheet';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import ErrorBoundary from 'Components/error-boundary';
import ListUsersTable from './subcomponents/list-users-table';

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
