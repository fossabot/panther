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
