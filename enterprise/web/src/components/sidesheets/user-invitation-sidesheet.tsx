/**
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program [The enterprise software] is licensed under the terms of a commercial license
 * available from Panther Labs Inc ("Panther Commercial License") by contacting contact@runpanther.com.
 * All use, distribution, and/or modification of this software, whether commercial or non-commercial,
 * falls under the Panther Commercial License to the extent it is permitted.
 */

import { Box, Heading, Text, SideSheet } from 'pouncejs';
import React from 'react';
import useSidesheet from 'Hooks/useSidesheet';
import InviteUserForm from '../../pages/users/subcomponents/invite-user-form';

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
