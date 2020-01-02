import React from 'react';
import { ComplianceItem } from 'Generated/schema';
import { Table, TableProps } from 'pouncejs';
import urls from 'Source/urls';
import useRouter from 'Hooks/useRouter';
import { generateEnumerationColumn } from 'Helpers/utils';

interface ResourcesDetailsTableProps {
  items?: ComplianceItem[];
  columns: TableProps<ComplianceItem>['columns'];
  enumerationStartIndex: number;
}

const ResourcesDetailsTable: React.FC<ResourcesDetailsTableProps> = ({
  enumerationStartIndex,
  items,
  columns,
}) => {
  const { history } = useRouter();

  // prepend an extra enumeration column
  const enumeratedColumns = [generateEnumerationColumn(enumerationStartIndex), ...columns];

  return (
    <Table<ComplianceItem>
      columns={enumeratedColumns}
      getItemKey={complianceItem => complianceItem.policyId}
      items={items}
      onSelect={complianceItem => history.push(urls.policies.details(complianceItem.policyId))}
    />
  );
};

export default React.memo(ResourcesDetailsTable);
