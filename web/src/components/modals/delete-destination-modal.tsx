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
