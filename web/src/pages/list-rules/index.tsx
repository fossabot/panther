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
import { Alert, Box, Card } from 'pouncejs';
import { DEFAULT_LARGE_PAGE_SIZE } from 'Source/constants';
import { useQuery, gql } from '@apollo/client';
import { convertObjArrayValuesToCsv, extractErrorMessage } from 'Helpers/utils';
import {
  ListRulesInput,
  ListRulesResponse,
  SortDirEnum,
  ListRulesSortFieldsEnum,
} from 'Generated/schema';

import TablePaginationControls from 'Components/utils/table-pagination-controls';
import useRequestParamsWithPagination from 'Hooks/useRequestParamsWithPagination';
import isEmpty from 'lodash-es/isEmpty';
import ErrorBoundary from 'Components/error-boundary';
import ListRulesTable from './subcomponents/list-rules-table';
import ListRulesActions from './subcomponents/list-rules-actions';
import ListRulesPageSkeleton from './skeleton';
import ListRulesPageEmptyDataFallback from './empty-data-fallback';

export const LIST_RULES = gql`
  query ListRules($input: ListRulesInput) {
    rules(input: $input) {
      rules {
        lastModified
        logTypes
        severity
        id
        displayName
        enabled
      }
      paging {
        totalPages
        thisPage
        totalItems
      }
    }
  }
`;

interface ApolloData {
  rules: ListRulesResponse;
}
interface ApolloVariables {
  input: ListRulesInput;
}

const ListRules = () => {
  const {
    requestParams,
    updateRequestParamsAndResetPaging,
    updatePagingParams,
  } = useRequestParamsWithPagination<ListRulesInput>();

  const { loading, error, data } = useQuery<ApolloData, ApolloVariables>(LIST_RULES, {
    fetchPolicy: 'cache-and-network',
    variables: {
      input: convertObjArrayValuesToCsv(requestParams),
    },
  });

  if (loading && !data) {
    return <ListRulesPageSkeleton />;
  }

  if (error) {
    return (
      <Alert
        mb={6}
        variant="error"
        title="Couldn't load your rules"
        description={
          extractErrorMessage(error) ||
          'There was an error when performing your request, please contact support@runpanther.io'
        }
      />
    );
  }

  // Get query results while protecting against exceptions
  const ruleItems = data.rules.rules;
  const pagingData = data.rules.paging;

  if (!ruleItems.length && isEmpty(requestParams)) {
    return <ListRulesPageEmptyDataFallback />;
  }

  //  Check how many active filters exist by checking how many columns keys exist in the URL
  return (
    <React.Fragment>
      <ListRulesActions />
      <ErrorBoundary>
        <Card>
          <ListRulesTable
            enumerationStartIndex={
              pagingData ? (pagingData.thisPage - 1) * DEFAULT_LARGE_PAGE_SIZE : 0
            }
            items={ruleItems}
            onSort={updateRequestParamsAndResetPaging}
            sortBy={requestParams.sortBy || ListRulesSortFieldsEnum.Id}
            sortDir={requestParams.sortDir || SortDirEnum.Ascending}
          />
        </Card>
      </ErrorBoundary>
      <Box my={5}>
        <TablePaginationControls
          page={pagingData.thisPage}
          totalPages={pagingData.totalPages}
          onPageChange={updatePagingParams}
        />
      </Box>
    </React.Fragment>
  );
};

export default ListRules;
