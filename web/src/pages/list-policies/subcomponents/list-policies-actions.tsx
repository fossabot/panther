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
import { RESOURCE_TYPES } from 'Source/constants';
import { ComplianceStatusEnum, SeverityEnum, ListPoliciesInput } from 'Generated/schema';
import GenerateFiltersGroup from 'Components/utils/generate-filters-group';
import { capitalize } from 'Helpers/utils';
import FormikTextInput from 'Components/fields/text-input';
import FormikCombobox from 'Components/fields/combobox';
import FormikMultiCombobox from 'Components/fields/multi-combobox';
import useRequestParamsWithPagination from 'Hooks/useRequestParamsWithPagination';
import { Box, Button, Card, Flex, Icon } from 'pouncejs';
import CreateButton from 'Pages/list-policies/subcomponents/create-button';
import ErrorBoundary from 'Components/error-boundary';
import isEmpty from 'lodash-es/isEmpty';
import pick from 'lodash-es/pick';

const severityOptions = Object.values(SeverityEnum);
const statusOptions = Object.values(ComplianceStatusEnum);

export const filters = {
  nameContains: {
    component: FormikTextInput,
    props: {
      label: 'Name contains',
      placeholder: 'Enter a policy name...',
    },
  },
  resourceTypes: {
    component: FormikMultiCombobox,
    props: {
      searchable: true,
      items: RESOURCE_TYPES,
      label: 'Resource Types',
      inputProps: {
        placeholder: 'Start typing resources...',
      },
    },
  },
  severity: {
    component: FormikCombobox,
    props: {
      label: 'Severity',
      items: severityOptions,
      itemToString: (severity: SeverityEnum) => capitalize(severity.toLowerCase()),
      inputProps: {
        placeholder: 'Choose a severity...',
      },
    },
  },
  tags: {
    component: FormikMultiCombobox,
    props: {
      label: 'Tags',
      searchable: true,
      allowAdditions: true,
      items: [] as string[],
      inputProps: {
        placeholder: 'Filter with tags...',
      },
    },
  },
  complianceStatus: {
    component: FormikCombobox,
    props: {
      label: 'Status',
      items: statusOptions,
      itemToString: (status: ComplianceStatusEnum) => capitalize(status.toLowerCase()),
      inputProps: {
        placeholder: 'Choose a status...',
      },
    },
  },
  hasRemediation: {
    component: FormikCombobox,
    props: {
      label: 'Auto-remediation Status',
      items: [true, false],
      itemToString: (item: boolean) => (item ? 'Configured' : 'Not Configured'),
      inputProps: {
        placeholder: 'Choose a remediation status...',
      },
    },
  },
};

export type ListPoliciesFiltersValues = Pick<
  ListPoliciesInput,
  'complianceStatus' | 'tags' | 'severity' | 'resourceTypes' | 'nameContains'
>;

const ListPoliciesActions: React.FC = () => {
  const [areFiltersVisible, setFiltersVisibility] = React.useState(false);
  const { requestParams, updateRequestParamsAndResetPaging } = useRequestParamsWithPagination<
    ListPoliciesInput
  >();

  const filterKeys = Object.keys(filters) as (keyof ListPoliciesInput)[];
  const filtersCount = filterKeys.filter(key => !isEmpty(requestParams[key])).length;

  // The initial filter values for when the filters component first renders. If you see down below,
  // we mount and unmount it depending on whether it's visible or not
  const initialFilterValues = React.useMemo(
    () => pick(requestParams, filterKeys) as ListPoliciesFiltersValues,
    [requestParams]
  );

  return (
    <Box>
      <Flex justifyContent="flex-end" mb={6}>
        <Box position="relative" mr={5}>
          <Button
            size="large"
            variant="default"
            onClick={() => setFiltersVisibility(!areFiltersVisible)}
          >
            <Flex>
              <Icon type="filter" size="small" mr={3} />
              Filter Options {filtersCount ? `(${filtersCount})` : ''}
            </Flex>
          </Button>
        </Box>
        <CreateButton />
      </Flex>
      {areFiltersVisible && (
        <ErrorBoundary>
          <Card p={6} mb={6}>
            <GenerateFiltersGroup<ListPoliciesFiltersValues>
              filters={filters}
              onCancel={() => setFiltersVisibility(false)}
              onSubmit={updateRequestParamsAndResetPaging}
              initialValues={initialFilterValues}
            />
          </Card>
        </ErrorBoundary>
      )}
    </Box>
  );
};

export default React.memo(ListPoliciesActions);
