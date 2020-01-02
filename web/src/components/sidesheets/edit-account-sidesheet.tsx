import { Box, Heading, SideSheet } from 'pouncejs';
import EditProfileForm from 'Components/forms/edit-profile-form';
import ChangePasswordForm from 'Components/forms/change-password-form';
import React from 'react';
import useSidesheet from 'Hooks/useSidesheet';

const EditAccountSidesheet: React.FC = () => {
  const { hideSidesheet } = useSidesheet();
  return (
    <SideSheet open onClose={hideSidesheet}>
      <Box mx={10} mb={10}>
        <Heading pt={1} pb={8} size="medium">
          Edit Profile
        </Heading>
        <EditProfileForm onSuccess={hideSidesheet} />
      </Box>
      <Box borderTop="1px solid" borderColor="grey100" mx={10}>
        <Heading py={8} size="medium">
          Account Security
        </Heading>
        <ChangePasswordForm />
      </Box>
    </SideSheet>
  );
};

export default EditAccountSidesheet;
