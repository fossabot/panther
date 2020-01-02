/* eslint-disable react/display-name */

import React from 'react';
import { Text, TableProps, Box } from 'pouncejs';
import { User } from 'Generated/schema';
import ListUsersTableRowOptions from 'Pages/users/subcomponents/list-users-table-row-options';
import dayjs from 'dayjs';
import { generateEnumerationColumn } from 'Helpers/utils';

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
