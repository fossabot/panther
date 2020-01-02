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
