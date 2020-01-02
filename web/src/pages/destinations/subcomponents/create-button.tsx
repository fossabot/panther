import React from 'react';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import { Button, Flex, Icon } from 'pouncejs';
import useSidesheet from 'Hooks/useSidesheet';

const DestinationCreateButton: React.FC = () => {
  const { showSidesheet } = useSidesheet();
  return (
    <Button
      size="large"
      variant="primary"
      mb={6}
      onClick={() =>
        showSidesheet({
          sidesheet: SIDESHEETS.SELECT_DESTINATION,
        })
      }
    >
      <Flex alignItems="center">
        <Icon type="add" size="small" mr={1} />
        Add Destination
      </Flex>
    </Button>
  );
};

export default DestinationCreateButton;
