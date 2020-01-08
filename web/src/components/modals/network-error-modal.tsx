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
import { Box, Flex, Modal, Text } from 'pouncejs';
import useModal from 'Hooks/useModal';
import SubmitButton from 'Components/utils/SubmitButton';

const NetworkErrorModal: React.FC = () => {
  const { hideModal } = useModal();
  return (
    <Modal open onClose={hideModal} title="No Internet Connection">
      <Box width={600}>
        <Text size="large" color="grey300" my={10} textAlign="center">
          Somebody is watching cat videos and is preventing you from being online
          <br />
          <br />
          That{"'"}s the most likely scenario anyway...
        </Text>
        <Flex justifyContent="center" mb={5}>
          <SubmitButton submitting disabled>
            Reconnecting
          </SubmitButton>
        </Flex>
      </Box>
    </Modal>
  );
};

export default NetworkErrorModal;
