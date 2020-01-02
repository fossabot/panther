/* eslint-disable react/display-name */

import React from 'react';
import { Badge, TableProps, Text } from 'pouncejs';
import { AlertSummary } from 'Generated/schema';
import { formatDatetime } from 'Helpers/utils';
import urls from 'Source/urls';
import { Link } from 'react-router-dom';
import { SEVERITY_COLOR_MAP } from 'Source/constants';

// The columns that the associated table will show
const columns = [
  {
    key: 'id',
    sortable: true,
    header: 'Alert ID',
    flex: '0 0 350px',
    renderCell: item => (
      <Link to={urls.alerts.details(item.alertId)}>
        <Text size="medium">{item.alertId}</Text>
      </Link>
    ),
  },
  {
    key: 'eventsMatched',
    sortable: true,
    header: 'Events Count',
    flex: '1 0 50px',
  },

  // Render badges to showcase severity
  {
    key: 'severity',
    sortable: true,
    flex: '1 0 100px',
    header: 'Severity',
    renderCell: item => <Badge color={SEVERITY_COLOR_MAP[item.severity]}>{item.severity}</Badge>,
  },

  // Date needs to be formatted properly
  {
    key: 'createdAt',
    sortable: true,
    header: 'Created At',
    flex: '0 0 200px',
    renderCell: ({ creationTime }) => <Text size="medium">{formatDatetime(creationTime)}</Text>,
  },
  // Date needs to be formatted properly
  {
    key: 'lastModified',
    sortable: true,
    header: 'Last Matched At',
    flex: '0 0 200px',
    renderCell: ({ lastEventMatched }) => (
      <Text size="medium">{formatDatetime(lastEventMatched)}</Text>
    ),
  },
] as TableProps<AlertSummary>['columns'];

export default columns;
