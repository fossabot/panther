/* eslint-disable react/display-name */

import React from 'react';
import { Text, TableProps } from 'pouncejs';
import { AlertSummary } from 'Generated/schema';
import { formatDatetime } from 'Helpers/utils';

// The columns that the associated table will show
const columns = [
  // The name is the `id` of the alert
  {
    key: 'alertId',
    header: 'Alert',
    flex: '2 0 450px',
  },

  {
    key: 'creationTime',
    header: 'Created At',
    flex: '1 0 250px',
    renderCell: ({ creationTime }) => <Text size="medium">{formatDatetime(creationTime)}</Text>,
  },
] as TableProps<AlertSummary>['columns'];

export default columns;
