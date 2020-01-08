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
import { Modal, Text, Flex, Button, useSnackbar } from 'pouncejs';
import { MutationTuple } from '@apollo/client';
import SubmitButton from 'Components/utils/SubmitButton';
import useModal from 'Hooks/useModal';

export interface BaseDeleteModalProps {
  mutation: MutationTuple<boolean, { [key: string]: any }>;
  itemDisplayName: string;
  onSuccess?: () => void;
  onError?: () => void;
}

const BaseDeleteModal: React.FC<BaseDeleteModalProps> = ({
  mutation,
  itemDisplayName,
  onSuccess = () => {},
  onError = () => {},
}) => {
  const { pushSnackbar } = useSnackbar();
  const { hideModal } = useModal();
  const [deleteItem, { loading, data, error }] = mutation;

  React.useEffect(() => {
    if (error) {
      pushSnackbar({ variant: 'error', title: `Failed to delete ${itemDisplayName}` });
      onError();
    }
  }, [error]);

  React.useEffect(() => {
    if (data) {
      pushSnackbar({ variant: 'success', title: `Successfully deleted ${itemDisplayName}` });
      hideModal();
      onSuccess();
    }
  }, [data]);

  return (
    <Modal open onClose={hideModal} title={`Delete ${itemDisplayName}`}>
      <Text size="large" color="grey500" mb={8} textAlign="center">
        Are you sure you want to delete <b>{itemDisplayName}</b>?
      </Text>

      <Flex justifyContent="flex-end">
        <Button size="large" variant="default" onClick={hideModal} mr={3}>
          Cancel
        </Button>
        <SubmitButton onClick={() => deleteItem()} submitting={loading} disabled={loading}>
          Delete
        </SubmitButton>
      </Flex>
    </Modal>
  );
};

export default BaseDeleteModal;
