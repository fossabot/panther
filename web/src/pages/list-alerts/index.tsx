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
import { Alert, Box, Card, Spinner } from 'pouncejs';
import { DEFAULT_LARGE_PAGE_SIZE } from 'Source/constants';
import { useQuery, gql } from '@apollo/client';
import { extractErrorMessage } from 'Helpers/utils';
import useInfiniteScroll from 'Hooks/useInfiniteScroll';
import { ListAlertsInput, ListAlertsResponse } from 'Generated/schema';
import ErrorBoundary from 'Components/error-boundary';
import ListAlertsTable from './subcomponents/list-alerts-table';
import ListAlertsPageSkeleton from './skeleton';
import ListAlertsPageEmptyDataFallback from './empty-data-fallback';

export const LIST_ALERTS = gql`
  query ListAlerts($input: ListAlertsInput) {
    alerts(input: $input) {
      alertSummaries {
        alertId
        creationTime
        eventsMatched
        lastEventMatched
        ruleId
        severity
      }
      lastEvaluatedKey
    }
  }
`;

interface ApolloData {
  alerts: ListAlertsResponse;
}
interface ApolloVariables {
  input: ListAlertsInput;
}

const ListAlerts = () => {
  const { loading, error, data, fetchMore } = useQuery<ApolloData, ApolloVariables>(LIST_ALERTS, {
    notifyOnNetworkStatusChange: true, // Adding notifyOnNetworkStatusChange will enable 'loading' to update its status during fetchMore requests as well
    fetchPolicy: 'cache-and-network',
    variables: {
      input: {
        pageSize: DEFAULT_LARGE_PAGE_SIZE,
      },
    },
  });

  const alertItems = data?.alerts.alertSummaries || [];
  const lastEvaluatedKey = data?.alerts.lastEvaluatedKey || null;
  const [infiniteRef, setHasNextPage] = useInfiniteScroll({
    loading,
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    onLoadMore: () => {
      // Even though we're setting hasNextPage as false when exclusiveStartKey is null
      // the react-infinite-scroll-hook library still makes one last request before finally stopping
      // We're adding this redundant check explicitly just to be sure
      if (!lastEvaluatedKey) {
        return;
      }

      fetchMore({
        variables: {
          input: { pageSize: DEFAULT_LARGE_PAGE_SIZE, exclusiveStartKey: lastEvaluatedKey },
        },
        updateQuery: (previousResult, { fetchMoreResult }: { fetchMoreResult: ApolloData }) => {
          if (!fetchMoreResult) {
            return previousResult;
          }
          const newAlertSummaries = fetchMoreResult.alerts.alertSummaries;
          const newLastEvaluatedKey = fetchMoreResult.alerts.lastEvaluatedKey;
          if (!newLastEvaluatedKey) {
            setHasNextPage(false); // newLastEvaluatedKey being null means there are no more items to query
          }
          return {
            alerts: {
              ...previousResult.alerts,
              alertSummaries: [...previousResult.alerts.alertSummaries, ...newAlertSummaries],
              lastEvaluatedKey: newLastEvaluatedKey,
            },
          };
        },
      });
    },
  });

  if (loading && !data) {
    return <ListAlertsPageSkeleton />;
  }

  if (error) {
    return (
      <Alert
        mb={6}
        variant="error"
        title="Couldn't load your alerts"
        description={
          extractErrorMessage(error) ||
          'There was an error when performing your request, please contact support@runpanther.io'
        }
      />
    );
  }

  if (!alertItems.length) {
    return <ListAlertsPageEmptyDataFallback />;
  }

  //  Check how many active filters exist by checking how many columns keys exist in the URL
  return (
    <ErrorBoundary>
      <div ref={infiniteRef}>
        <Card mb={8}>
          <ListAlertsTable items={alertItems} />
        </Card>
        {loading && (
          <Box mb={8}>
            <Spinner size="large" margin="auto" />
          </Box>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default ListAlerts;
