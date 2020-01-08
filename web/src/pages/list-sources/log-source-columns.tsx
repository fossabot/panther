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
import { TableProps, Box } from 'pouncejs';
import { Integration } from 'Generated/schema';
import { generateEnumerationColumn } from 'Helpers/utils';
import ListSourcesTableRowOptionsProps from 'Pages/list-sources/subcomponents/list-sources-table-row-options';

// The columns that the associated table will show
const columns = [
  generateEnumerationColumn(0),

  // The source label that user defined
  {
    key: 'integrationLabel',
    header: 'Label',
    flex: '1 0 200px',
  },

  {
    key: 'sourceSnsTopicArn',
    header: 'SNS Topic ARN',
    flex: '1 0 200px',
  },

  {
    key: 'logProcessingRoleArn',
    header: 'Assumable Role ARN',
    flex: '1 0 200px',
  },

  {
    key: 'options',
    flex: '0 1 auto',
    renderColumnHeader: () => <Box mx={5} />,
    renderCell: item => <ListSourcesTableRowOptionsProps source={item} />,
  },
] as TableProps<Integration>['columns'];

export default columns;
