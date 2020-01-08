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
import { Text, TableProps, Badge, Tooltip, Label } from 'pouncejs';
import { ComplianceItem, ComplianceStatusEnum } from 'Generated/schema';
import { capitalize } from 'Helpers/utils';
import { SEVERITY_COLOR_MAP } from 'Source/constants';
import RemediationButton from 'Components/table-row-remediation-button';
import SuppressButton from 'Components/table-row-suppress-button';

// The columns that the associated table will show
const columns = [
  // The name is the `id` of the resource
  {
    key: 'policyId',
    header: 'Policy',
    flex: '1 0 450px',
  },

  // Showcase the status (pass/fail) with the proper text color
  {
    key: 'status',
    header: 'Status',
    flex: '0 0 125px',
    renderCell: ({ status, errorMessage }) => {
      const hasErrored = status === ComplianceStatusEnum.Error;
      const textNode = (
        <Text size="medium" color={status === ComplianceStatusEnum.Pass ? 'green300' : 'red300'}>
          {capitalize(status.toLowerCase())}
          {hasErrored && ' *'}
        </Text>
      );

      if (errorMessage) {
        return (
          <Tooltip positioning="down" content={<Label size="medium">{errorMessage}</Label>}>
            {textNode}
          </Tooltip>
        );
      }

      return textNode;
    },
  },

  // Render badges to showcase severity
  {
    key: 'severity',
    flex: '0 0 125px',
    header: 'Severity',
    renderCell: ({ policySeverity }) => (
      <Badge color={SEVERITY_COLOR_MAP[policySeverity]}>{policySeverity}</Badge>
    ),
  },

  // The remediation button
  {
    key: 'remediationOptions',
    flex: '0 0 140px',
    renderColumnHeader: () => null,
    renderCell: (complianceItem, index) => {
      return (
        complianceItem.status !== ComplianceStatusEnum.Pass && (
          <RemediationButton
            buttonVariant={index % 2 === 0 ? 'default' : 'secondary'}
            policyId={complianceItem.policyId}
            resourceId={complianceItem.resourceId}
          />
        )
      );
    },
  },

  // The suppress button
  {
    key: 'suppressionOptions',
    flex: '0 0 120px',
    renderColumnHeader: () => null,
    renderCell: (complianceItem, index) => {
      return !complianceItem.suppressed ? (
        <SuppressButton
          buttonVariant={index % 2 === 0 ? 'default' : 'secondary'}
          policyIds={[complianceItem.policyId]}
          resourcePatterns={[complianceItem.resourceId]}
        />
      ) : (
        <Label color="orange300" size="medium">
          IGNORED
        </Label>
      );
    },
  },
] as TableProps<ComplianceItem>['columns'];

export default columns;
