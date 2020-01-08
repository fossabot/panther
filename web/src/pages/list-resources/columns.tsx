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
import { Text, TableProps, Tooltip, Label } from 'pouncejs';
import { ComplianceStatusEnum, Integration, ResourceSummary } from 'Generated/schema';
import { capitalize, formatDatetime } from 'Helpers/utils';

// The columns that the associated table will show
const columns = [
  // The name is the `id` of the resource
  {
    key: 'id',
    sortable: true,
    header: 'Resource',
    flex: '0 0 350px',
  },

  // The AWS type of this resouce (S3, IAM, etc.)
  {
    key: 'type',
    sortable: true,
    header: 'Type',
    flex: '0 1 275px',
  },

  // The AWS account associated with this resource within the context of an organization
  {
    key: 'integrationLabel',
    sortable: true,
    header: 'Source',
    flex: '1 0 100px',
  },

  // Status is not available yet. Mock it by alternative between hardcoded text
  {
    key: 'complianceStatus',
    sortable: true,
    header: 'Status',
    flex: '0 0 100px',
    renderCell: ({ complianceStatus }) => {
      const hasErrored = complianceStatus === ComplianceStatusEnum.Error;
      const textNode = (
        <Text
          size="medium"
          color={complianceStatus === ComplianceStatusEnum.Pass ? 'green300' : 'red300'}
        >
          {capitalize(complianceStatus.toLowerCase())}
          {hasErrored && ' *'}
        </Text>
      );

      if (hasErrored) {
        return (
          <Tooltip
            positioning="down"
            content={
              <Label size="medium">
                Some policies have raised an exception when evaluating this resource. Find out more
                in the resource{"'"}s page
              </Label>
            }
          >
            {textNode}
          </Tooltip>
        );
      }

      return textNode;
    },
  },

  // Date needs to be formatted properly
  {
    key: 'lastModified',
    sortable: true,
    header: 'Last Modified',
    flex: '0 1 225px',
    renderCell: ({ lastModified }) => <Text size="medium">{formatDatetime(lastModified)}</Text>,
  },
] as TableProps<ResourceSummary & Pick<Integration, 'integrationLabel'>>['columns'];

export default columns;
