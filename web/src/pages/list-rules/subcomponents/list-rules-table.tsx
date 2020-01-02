import React from 'react';
import {
  ListRulesInput,
  ListRulesSortFieldsEnum,
  RuleSummary,
  SortDirEnum,
} from 'Generated/schema';
import { generateEnumerationColumn } from 'Helpers/utils';
import { Table } from 'pouncejs';
import columns from 'Pages/list-rules/columns';
import urls from 'Source/urls';
import useRouter from 'Hooks/useRouter';

interface ListRulesTableProps {
  items?: RuleSummary[];
  sortBy: ListRulesSortFieldsEnum;
  sortDir: SortDirEnum;
  onSort: (params: Partial<ListRulesInput>) => void;
  enumerationStartIndex: number;
}

const ListRulesTable: React.FC<ListRulesTableProps> = ({
  items,
  onSort,
  sortBy,
  sortDir,
  enumerationStartIndex,
}) => {
  const { history } = useRouter();

  const handleSort = (selectedKey: ListRulesSortFieldsEnum) => {
    if (sortBy === selectedKey) {
      onSort({
        sortBy,
        sortDir: sortDir === SortDirEnum.Ascending ? SortDirEnum.Descending : SortDirEnum.Ascending,
      });
    } else {
      onSort({ sortBy: selectedKey, sortDir: SortDirEnum.Ascending });
    }
  };

  const enumeratedColumns = [generateEnumerationColumn(enumerationStartIndex), ...columns];

  return (
    <Table<RuleSummary>
      columns={enumeratedColumns}
      getItemKey={rule => rule.id}
      items={items}
      onSort={handleSort}
      sortDir={sortDir}
      sortKey={sortBy}
      onSelect={rule => history.push(urls.rules.details(rule.id))}
    />
  );
};

export default React.memo(ListRulesTable);
