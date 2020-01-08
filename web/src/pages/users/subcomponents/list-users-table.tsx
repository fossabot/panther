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

import React from 'react';
import { useQuery, gql } from '@apollo/client';
import { ListUsersOutput, User } from 'Generated/schema';
import { Alert, Card, Table } from 'pouncejs';
import columns from 'Pages/users/columns';

import TablePlaceholder from 'Components/table-placeholder';
import { extractErrorMessage } from 'Helpers/utils';

// This is done so we can benefit from React.memo
const getUserItemKey = (item: User) => item.id;

export const LIST_USERS = gql`
  query ListUsers($limit: Int, $paginationToken: String) {
    users(limit: $limit, paginationToken: $paginationToken) {
      users {
        id
        email
        givenName
        familyName
        createdAt
        status
        role
      }
      paginationToken
    }
  }
`;

const ListUsersTable = () => {
  const { loading, error, data } = useQuery<{ users: ListUsersOutput }>(LIST_USERS, {
    fetchPolicy: 'cache-and-network',
  });

  if (loading && !data) {
    return (
      <Card p={9}>
        <TablePlaceholder />
      </Card>
    );
  }

  if (error) {
    return (
      <Alert
        variant="error"
        title="Couldn't load users"
        description={
          extractErrorMessage(error) ||
          'There was an error when performing your request, please contact support@runpanther.io'
        }
      />
    );
  }

  return (
    <Card>
      <Table<User> columns={columns} getItemKey={getUserItemKey} items={data.users.users} />
    </Card>
  );
};

export default React.memo(ListUsersTable);
