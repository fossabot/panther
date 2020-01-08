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
import ListInfraSourcesTable from 'Pages/list-sources/subcomponents/infra-source-table';
import ListLogSourcesTable from 'Pages/list-sources/subcomponents/log-source-table';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import { INTEGRATION_TYPES, READONLY_ROLES_ARRAY } from 'Source/constants';
import { Link } from 'react-router-dom';
import urls from 'Source/urls';
import ErrorBoundary from 'Components/error-boundary';
import Panel from 'Components/panel';

const ListSources = () => {
  return (
    <Box mb={6}>
      <Box mb={6}>
        <Panel
          title="AWS Account Sources"
          size="large"
          actions={
            <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
              <Button
                size="large"
                variant="primary"
                is={Link}
                to={urls.account.settings.sources.create(INTEGRATION_TYPES.AWS_INFRA)}
              >
                <Flex alignItems="center">
                  <Icon type="add" size="small" mr={1} />
                  Add Account
                </Flex>
              </Button>
            </RoleRestrictedAccess>
          }
        >
          <ErrorBoundary>
            <ListInfraSourcesTable />
          </ErrorBoundary>
        </Panel>
      </Box>
      <Box>
        <Panel
          title="Log Sources"
          size="large"
          actions={
            <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
              <Button
                size="large"
                variant="primary"
                is={Link}
                to={urls.account.settings.sources.create(INTEGRATION_TYPES.AWS_LOGS)}
              >
                <Flex alignItems="center">
                  <Icon type="add" size="small" mr={1} />
                  Add Source
                </Flex>
              </Button>
            </RoleRestrictedAccess>
          }
        >
          <ErrorBoundary>
            <ListLogSourcesTable />
          </ErrorBoundary>
        </Panel>
      </Box>
    </Box>
  );
};

export default ListSources;
