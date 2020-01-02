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
