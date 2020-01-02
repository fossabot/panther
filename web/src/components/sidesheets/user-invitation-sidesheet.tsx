import { Box, Heading, Text, SideSheet } from 'pouncejs';
import InviteUserForm from 'Pages/users/subcomponents/invite-user-form';
import React from 'react';
import useSidesheet from 'Hooks/useSidesheet';

const UserInvitationSidesheet: React.FC = () => {
  const { hideSidesheet } = useSidesheet();

  return (
    <SideSheet open onClose={hideSidesheet}>
      <Box width={460}>
        <Heading size="medium" mb={8}>
          Invite User
        </Heading>
        <Text size="large" color="grey200" mb={8}>
          By inviting users to join your organization, they will receive an email with temporary
          credentials that they can use to sign in to the platform
        </Text>
        <InviteUserForm onSuccess={hideSidesheet} />
      </Box>
    </SideSheet>
  );
};

export default UserInvitationSidesheet;
