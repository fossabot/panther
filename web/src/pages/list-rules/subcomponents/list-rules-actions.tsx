import React from 'react';
import { SeverityEnum, ListRulesInput } from 'Generated/schema';
import GenerateFiltersGroup from 'Components/utils/generate-filters-group';
import { capitalize } from 'Helpers/utils';
import { LOG_TYPES } from 'Source/constants';
import FormikCombobox from 'Components/fields/combobox';
import FormikMultiCombobox from 'Components/fields/multi-combobox';
import FormikTextInput from 'Components/fields/text-input';
import { Box, Button, Card, Flex, Icon } from 'pouncejs';
import CreateButton from 'Pages/list-rules/subcomponents/create-button';
import ErrorBoundary from 'Components/error-boundary';
import useRequestParamsWithPagination from 'Hooks/useRequestParamsWithPagination';
import isEmpty from 'lodash-es/isEmpty';
import pick from 'lodash-es/pick';

const severityOptions = Object.values(SeverityEnum);

export const filters = {
  nameContains: {
    component: FormikTextInput,
    props: {
      label: 'Name contains',
      placeholder: 'Enter a rule name...',
    },
  },
  logTypes: {
    component: FormikMultiCombobox,
    props: {
      searchable: true,
      items: LOG_TYPES,
      label: 'Log Types',
      inputProps: {
        placeholder: 'Start typing logs...',
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
};

export type ListRulesFiltersValues = Pick<
  ListRulesInput,
  'tags' | 'severity' | 'logTypes' | 'nameContains'
>;

const ListRulesActions: React.FC = () => {
  const [areFiltersVisible, setFiltersVisibility] = React.useState(false);
  const { requestParams, updateRequestParamsAndResetPaging } = useRequestParamsWithPagination<
    ListRulesInput
  >();

  const filterKeys = Object.keys(filters) as (keyof ListRulesInput)[];
  const filtersCount = filterKeys.filter(key => !isEmpty(requestParams[key])).length;

  // The initial filter values for when the filters component first renders. If you see down below,
  // we mount and unmount it depending on whether it's visible or not
  const initialFilterValues = React.useMemo(
    () => pick(requestParams, filterKeys) as ListRulesFiltersValues,
    [requestParams]
  );

  return (
    <React.Fragment>
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
            <GenerateFiltersGroup<ListRulesFiltersValues>
              filters={filters}
              onCancel={() => setFiltersVisibility(false)}
              onSubmit={updateRequestParamsAndResetPaging}
              initialValues={initialFilterValues}
            />
          </Card>
        </ErrorBoundary>
      )}
    </React.Fragment>
  );
};

export default React.memo(ListRulesActions);
