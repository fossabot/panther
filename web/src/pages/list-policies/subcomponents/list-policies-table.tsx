import React from 'react';
import {
  ListPoliciesInput,
  ListPoliciesSortFieldsEnum,
  PolicySummary,
  SortDirEnum,
} from 'Generated/schema';
import { generateEnumerationColumn } from 'Helpers/utils';
import { Table } from 'pouncejs';
import columns from 'Pages/list-policies/columns';
import urls from 'Source/urls';
import useRouter from 'Hooks/useRouter';

interface ListPoliciesTableProps {
  items?: PolicySummary[];
  sortBy: ListPoliciesSortFieldsEnum;
  sortDir: SortDirEnum;
  onSort: (params: Partial<ListPoliciesInput>) => void;
  enumerationStartIndex: number;
}

const ListPoliciesTable: React.FC<ListPoliciesTableProps> = ({
  items,
  onSort,
  sortBy,
  sortDir,
  enumerationStartIndex,
}) => {
  const { history } = useRouter();

  const handleSort = (selectedKey: ListPoliciesSortFieldsEnum) => {
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
    <Table<PolicySummary>
      columns={enumeratedColumns}
      getItemKey={policy => policy.id}
      items={items}
      onSort={handleSort}
      sortDir={sortDir}
      sortKey={sortBy}
      onSelect={policy => history.push(urls.policies.details(policy.id))}
    />
  );
};

export default React.memo(ListPoliciesTable);
