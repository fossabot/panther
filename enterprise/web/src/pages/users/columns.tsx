/**
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program [The enterprise software] is licensed under the terms of a commercial license
 * available from Panther Labs Inc ("Panther Commercial License") by contacting contact@runpanther.com.
 * All use, distribution, and/or modification of this software, whether commercial or non-commercial,
 * falls under the Panther Commercial License to the extent it is permitted.
 */

/* eslint-disable react/display-name */

import React from 'react';
import { Text, TableProps, Box } from 'pouncejs';
import { User } from 'Generated/schema';
import dayjs from 'dayjs';
import { generateEnumerationColumn } from 'Helpers/utils';
import ListUsersTableRowOptions from './subcomponents/list-users-table-row-options';

// The columns that the associated table will show
const columns = [
  generateEnumerationColumn(0),
  // Show given name and family name in two separate column
  {
    key: 'givenName',
    header: 'First Name',
    flex: '0 0 10%',
  },
  {
    key: 'familyName',
    header: 'Last Name',
    flex: '0 0 10%',
  },
  {
    key: 'email',
    header: 'Email',
    flex: '0 0 22%',
  },
  // Display user roles Admin, Analyst or ReadOnly
  {
    key: 'role',
    header: 'Role',
    flex: '0 0 8%',
  },
  // Display when user is invited
  {
    key: 'createdAt',
    header: 'Invited at',
    flex: '0 0 18%',
    renderCell: item => (
      <Text size="medium">{dayjs(item.createdAt * 1000).format('MM/DD/YYYY, HH:mm G[M]TZZ')}</Text>
    ),
  },
  // Display if user is confirmed or not
  {
    key: 'status',
    header: 'Status',
    flex: '0 0 250px',
  },
  {
    key: 'options',
    flex: '0 1 auto',
    renderColumnHeader: () => <Box mx={5} />,
    renderCell: item => <ListUsersTableRowOptions user={item} />,
  },
] as TableProps<User>['columns'];

export default columns;
