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
import { Destination } from 'Generated/schema';
import { useMutation, gql } from '@apollo/client';
import { LIST_DESTINATIONS } from 'Pages/destinations';
import BaseDeleteModal from 'Components/modals/base-delete-modal';

const DELETE_DESTINATION = gql`
  mutation DeleteOutput($id: ID!) {
    deleteDestination(id: $id)
  }
`;

export interface DeleteDestinationModalProps {
  destination: Destination;
}

export interface ApolloMutationInput {
  id: string;
}

const DeleteDestinationModal: React.FC<DeleteDestinationModalProps> = ({ destination }) => {
  const destinationDisplayName = destination.displayName || destination.outputId;
  const mutation = useMutation<boolean, ApolloMutationInput>(DELETE_DESTINATION, {
    variables: {
      id: destination.outputId,
    },
    refetchQueries: [{ query: LIST_DESTINATIONS }],
  });

  return <BaseDeleteModal mutation={mutation} itemDisplayName={destinationDisplayName} />;
};

export default DeleteDestinationModal;
