import React from 'react';
import { useQuery, gql } from '@apollo/client';
import { ListUsersOutput, User } from 'Generated/schema';
import { Alert, Card, Table } from 'pouncejs';
import TablePlaceholder from 'Components/table-placeholder';
import { extractErrorMessage } from 'Helpers/utils';
import columns from '../columns';

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
