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
