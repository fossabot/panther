import React from 'react';
import { Alert, Box, Card, Flex, Table } from 'pouncejs';
import { READONLY_ROLES_ARRAY } from 'Source/constants';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import ErrorBoundary from 'Components/error-boundary';
import { gql, useQuery } from '@apollo/client';
import { Destination } from 'Generated/schema';
import { extractErrorMessage } from 'Helpers/utils';
import columns from './columns';
import DestinationsPageSkeleton from './skeleton';
import DestinationsPageEmptyDataFallback from './empty-data-fallback';
import DestinationCreateButton from './subcomponents/create-button';

export const LIST_DESTINATIONS = gql`
  query ListDestinationsAndDefaults {
    destinations {
      createdBy
      creationTime
      displayName
      lastModifiedBy
      lastModifiedTime
      outputId
      outputType
      outputConfig {
        slack {
          webhookURL
        }
        sns {
          topicArn
        }
        email {
          destinationAddress
        }
        pagerDuty {
          integrationKey
        }
        github {
          repoName
          token
        }
        jira {
          orgDomain
          projectKey
          userName
          apiKey
          assigneeID
        }
        opsgenie {
          apiKey
        }
        msTeams {
          webhookURL
        }
      }
      verificationStatus
      defaultForSeverity
    }
  }
`;

export interface ListDestinationsQueryData {
  destinations: Destination[];
}

const ListDestinations = () => {
  const { loading, error, data } = useQuery<ListDestinationsQueryData>(LIST_DESTINATIONS, {
    fetchPolicy: 'cache-and-network',
  });

  if (loading && !data) {
    return <DestinationsPageSkeleton />;
  }

  if (error) {
    return (
      <Alert
        variant="error"
        title="Couldn't load your available destinations"
        description={
          extractErrorMessage(error) ||
          'There was an error while attempting to list your destinations'
        }
      />
    );
  }

  if (!data.destinations.length) {
    return <DestinationsPageEmptyDataFallback />;
  }

  return (
    <Box mb={6}>
      <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
        <Flex justifyContent="flex-end">
          <DestinationCreateButton />
        </Flex>
      </RoleRestrictedAccess>
      <Card>
        <ErrorBoundary>
          <Table<Destination>
            items={data.destinations}
            getItemKey={item => item.outputId}
            columns={columns}
          />
        </ErrorBoundary>
      </Card>
    </Box>
  );
};

export default ListDestinations;
