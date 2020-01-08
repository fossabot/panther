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
import { ComplianceItem, Integration } from 'Generated/schema';
import { Table, TableProps } from 'pouncejs';
import urls from 'Source/urls';
import useRouter from 'Hooks/useRouter';
import { generateEnumerationColumn } from 'Helpers/utils';

type EnhancedComplianceItem = ComplianceItem & Pick<Integration, 'integrationLabel'>;

interface PolicyDetailsTableProps {
  items?: EnhancedComplianceItem[];
  columns: TableProps<EnhancedComplianceItem>['columns'];
  enumerationStartIndex: number;
}

const PolicyDetailsTable: React.FC<PolicyDetailsTableProps> = ({
  items,
  columns,
  enumerationStartIndex,
}) => {
  const { history } = useRouter();

  // prepend an extra enumeration column
  const enumeratedColumns = [generateEnumerationColumn(enumerationStartIndex), ...columns];

  return (
    <Table<EnhancedComplianceItem>
      columns={enumeratedColumns}
      getItemKey={complianceItem => complianceItem.resourceId}
      items={items}
      onSelect={complianceItem => history.push(urls.resources.details(complianceItem.resourceId))}
    />
  );
};

export default React.memo(PolicyDetailsTable);
