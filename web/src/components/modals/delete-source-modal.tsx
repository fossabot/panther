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
