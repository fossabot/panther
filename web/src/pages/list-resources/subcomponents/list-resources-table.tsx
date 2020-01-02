import React from 'react';
import {
  ListResourcesSortFieldsEnum,
  ListResourcesInput,
  ResourceSummary,
  SortDirEnum,
  Integration,
} from 'Generated/schema';
import { generateEnumerationColumn } from 'Helpers/utils';
import { Table } from 'pouncejs';
import useRouter from 'Hooks/useRouter';
import urls from 'Source/urls';
import columns from '../columns';

interface ListResourcesTableProps {
  items?: ResourceSummary[];
  sortBy: ListResourcesSortFieldsEnum;
  sortDir: SortDirEnum;
  onSort: (params: Partial<ListResourcesInput>) => void;
  enumerationStartIndex: number;
}

const ListResourcesTable: React.FC<ListResourcesTableProps> = ({
  items,
  onSort,
  sortBy,
  sortDir,
  enumerationStartIndex,
}) => {
  const { history } = useRouter();

  const handleSort = (selectedKey: ListResourcesSortFieldsEnum) => {
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
    <Table<ResourceSummary & Pick<Integration, 'integrationLabel'>>
      columns={enumeratedColumns}
      getItemKey={resource => resource.id}
      items={items}
      onSort={handleSort}
      sortDir={sortDir}
      sortKey={sortBy}
      onSelect={resource => history.push(urls.resources.details(resource.id))}
    />
  );
};

export default React.memo(ListResourcesTable);
