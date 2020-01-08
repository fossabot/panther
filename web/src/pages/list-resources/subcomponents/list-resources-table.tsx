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
