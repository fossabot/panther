/* eslint-disable react/display-name */

import React from 'react';
import { Text, TableProps, Icon, Box } from 'pouncejs';
import { Integration } from 'Generated/schema';
import { generateEnumerationColumn } from 'Helpers/utils';
import ListSourcesTableRowOptionsProps from 'Pages/list-sources/subcomponents/list-sources-table-row-options';

// The columns that the associated table will show
const columns = [
  generateEnumerationColumn(0),

  // The account is the `id` number of the aws account
  {
    key: 'awsAccountId',
    header: 'Account',
    flex: '1 0 150px',
  },

  // The source label that user defined
  {
    key: 'integrationLabel',
    header: 'Label',
    flex: '1 0 275px',
  },

  // Status displays the error message
  {
    key: 'lastScanErrorMessage',
    header: 'Status',
    flex: '1 0 150px',
    renderCell: item => {
      const isFailing = Boolean(item.lastScanErrorMessage);
      if (!isFailing) {
        return <Icon color="green300" size="small" type="check" />;
      }
      return (
        <Text size="medium" color="red300">
          {item.lastScanErrorMessage}
        </Text>
      );
    },
  },
  {
    key: 'options',
    flex: '0 1 auto',
    renderColumnHeader: () => <Box mx={5} />,
    renderCell: item => <ListSourcesTableRowOptionsProps source={item} />,
  },
] as TableProps<Integration>['columns'];

export default columns;
