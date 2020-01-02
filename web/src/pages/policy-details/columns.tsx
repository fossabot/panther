/* eslint-disable react/display-name */

import React from 'react';
import { Text, TableProps, Tooltip, Label } from 'pouncejs';
import { ComplianceItem, ComplianceStatusEnum } from 'Generated/schema';
import { capitalize, formatDatetime } from 'Helpers/utils';
import RemediationButton from 'Components/table-row-remediation-button';
import SuppressButton from 'Components/table-row-suppress-button';

// The columns that the associated table will show
const columns = [
  // The name is the `id` of the resource
  {
    key: 'resourceId',
    header: 'Resource',
    flex: '0 0 400px',
  },

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
  {
    key: 'integrationLabel',
    flex: '1 0 100px',
    header: 'Source',
  },

  // Render badges to showcase severity
  {
    key: 'lastUpdated',
    flex: '1 0 170px',
    header: 'Last Updated',
    renderCell: ({ lastUpdated }) => <Text size="medium">{formatDatetime(lastUpdated)}</Text>,
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
        <Label color="orange300" size="medium" mx={2}>
          IGNORED
        </Label>
      );
    },
  },
] as TableProps<ComplianceItem>['columns'];

export default columns;
