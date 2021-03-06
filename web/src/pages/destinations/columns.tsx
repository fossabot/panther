/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

/* eslint-disable react/display-name */

import React from 'react';
import { TableProps, Box, Badge, Text, Tooltip, Label } from 'pouncejs';
import { formatDatetime, generateEnumerationColumn } from 'Helpers/utils';
import ListDestinationsTableRowOptionsProps from 'Pages/destinations/subcomponents/list-destinations-table-row-options';
import { SEVERITY_COLOR_MAP } from 'Source/constants';
import { DestinationTypeEnum, Destination } from 'Generated/schema';

// The columns that the associated table will show
const columns = [
  generateEnumerationColumn(0),

  // The user specified display name for the destination
  {
    key: 'displayName',
    header: 'Display Name',
    flex: '1 0 200px',
    renderCell: ({ displayName, verificationStatus, outputType, outputConfig }) => {
      const isUnverifiedEmailDestination =
        outputType === DestinationTypeEnum.Email && verificationStatus !== 'SUCCESS';

      if (!isUnverifiedEmailDestination) {
        return <Text size="medium">{displayName}</Text>;
      }

      const emailAddress = outputConfig.email.destinationAddress;
      let verificationMessage;
      switch (verificationStatus) {
        case 'PENDING':
          verificationMessage = `${emailAddress} is currently pending verification`;
          break;
        case 'FAILED':
          verificationMessage = `${emailAddress} failed to become verified. Please update it`;
          break;
        default:
          verificationMessage = "Email verification process hasn't been initiated";
      }
      return (
        <Tooltip positioning="down" content={<Label size="medium">{verificationMessage}</Label>}>
          <Text size="medium" color="red300">
            {displayName} *
          </Text>
        </Tooltip>
      );
    },
  },

  // The service like slack or pagerduty
  {
    key: 'outputType',
    header: 'Integrated Service',
    flex: '1 0 175px',
  },

  // Default severities this destination is assigned to
  {
    key: 'defaultForSeverity',
    header: 'Associated Severities',
    flex: '0 1 375px',
    renderCell: (item: Destination) => {
      return item.defaultForSeverity.map(severity => (
        <Badge key={`${item.outputId}${severity}`} color={SEVERITY_COLOR_MAP[severity]} mr={1}>
          {severity}
        </Badge>
      ));
    },
  },

  // The time that it was created
  {
    key: 'creationTime',
    sortable: true,
    header: 'Created at',
    flex: '0 0 225px',
    renderCell: ({ creationTime }) => <Text size="medium">{formatDatetime(creationTime)}</Text>,
  },

  {
    key: 'options',
    flex: '0 1 auto',
    renderColumnHeader: () => <Box width={30} />,
    renderCell: (item: Destination) => <ListDestinationsTableRowOptionsProps destination={item} />,
  },
] as TableProps<Destination>['columns'];

export default columns;
