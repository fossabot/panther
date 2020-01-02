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
