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
import { generateEnumerationColumn } from 'Helpers/utils';
import { Badge, TableProps } from 'pouncejs';
import { SEVERITY_COLOR_MAP } from 'Source/constants';
import { TopFailingPolicy, TopFailingResource } from 'Pages/overview/index';

/**
 * The columns that the top failing policies table will show
 */
export const topFailingPoliciesColumns = [
  // add an enumeration column starting from 0
  generateEnumerationColumn(0),

  // The name is the `id` of the policy
  {
    key: 'id',
    header: 'Policy',
    flex: '1 0 0',
  },

  // Render badges to showcase severity
  {
    key: 'severity',
    flex: '0 0 100px',
    header: 'Severity',
    renderCell: ({ severity }) => <Badge color={SEVERITY_COLOR_MAP[severity]}>{severity}</Badge>,
  },
] as TableProps<TopFailingPolicy>['columns'];

/**
 * The columns that the top failing resources table will show
 */
export const topFailingResourcesColumns = [
  // add an enumeration column starting from 0
  generateEnumerationColumn(0),

  // The name is the `id` of the policy
  {
    key: 'id',
    header: 'Resource',
    flex: '1 0 0',
  },
] as TableProps<TopFailingResource>['columns'];
