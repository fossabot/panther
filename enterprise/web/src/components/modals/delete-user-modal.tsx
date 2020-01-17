/**
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program [The enterprise software] is licensed under the terms of a commercial license
 * available from Panther Labs Inc ("Panther Commercial License") by contacting contact@runpanther.com.
 * All use, distribution, and/or modification of this software, whether commercial or non-commercial,
 * falls under the Panther Commercial License to the extent it is permitted.
 */

import React from 'react';
import { User } from 'Generated/schema';

import { useMutation, gql } from '@apollo/client';
import { LIST_USERS } from 'Pages/users/subcomponents/list-users-table';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import BaseDeleteModal from 'Components/modals/base-delete-modal';

const DELETE_USER = gql`
  mutation DeleteUser($id: ID!) {
    removeUser(id: $id)
  }
`;

export interface DeleteUserModalProps {
  user: User;
}

const DeleteUserModal: React.FC<DeleteUserModalProps> = ({ user }) => {
  const userDisplayName = `${user.givenName} ${user.familyName}` || user.id;
  const mutation = useMutation<boolean, { id: string }>(DELETE_USER, {
    variables: {
      id: user.id,
    },
    awaitRefetchQueries: true,
    refetchQueries: [getOperationName(LIST_USERS)],
  });

  return <BaseDeleteModal mutation={mutation} itemDisplayName={userDisplayName} />;
};

export default DeleteUserModal;
