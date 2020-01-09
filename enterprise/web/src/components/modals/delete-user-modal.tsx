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
