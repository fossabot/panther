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
