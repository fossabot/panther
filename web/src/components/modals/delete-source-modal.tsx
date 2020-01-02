import React from 'react';
import { Integration } from 'Generated/schema';
import { LIST_INFRA_SOURCES } from 'Pages/list-sources/subcomponents/infra-source-table';
import { LIST_LOG_SOURCES } from 'Pages/list-sources/subcomponents/log-source-table';
import { useMutation, gql } from '@apollo/client';
import BaseDeleteModal from 'Components/modals/base-delete-modal';
import { INTEGRATION_TYPES } from 'Source/constants';

const DELETE_SOURCE = gql`
  mutation DeleteSource($id: ID!) {
    deleteIntegration(id: $id)
  }
`;

export interface DeleteSourceModalProps {
  source: Integration;
}

const DeleteSourceModal: React.FC<DeleteSourceModalProps> = ({ source }) => {
  const isInfraSource = source.integrationType === INTEGRATION_TYPES.AWS_INFRA;
  const sourceDisplayName = source.integrationLabel || source.integrationId;
  const mutation = useMutation<boolean, { id: string }>(DELETE_SOURCE, {
    variables: {
      id: source.integrationId,
    },
    refetchQueries: [{ query: isInfraSource ? LIST_INFRA_SOURCES : LIST_LOG_SOURCES }],
  });

  return <BaseDeleteModal mutation={mutation} itemDisplayName={sourceDisplayName} />;
};

export default DeleteSourceModal;
