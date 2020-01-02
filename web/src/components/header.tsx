import React from 'react';
import Breadcrumbs from 'Components/breadcrumbs';
import { Button, Flex, Icon, IconButton, Text, Dropdown, MenuItem } from 'pouncejs';
import useAuth from 'Hooks/useAuth';
import useSidesheet from 'Hooks/useSidesheet';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';

const Header = () => {
  const { userInfo, signOut } = useAuth();
  const { showSidesheet } = useSidesheet();

  const userButton = React.useMemo(
    () => (
      <Button size="small" variant="default" my="auto" is="div">
        <Flex alignItems="center">
          <Icon type="user" size="small" mr={2} borderRadius="circle" bg="grey200" color="white" />
          {userInfo && (
            <Text size="medium">
              {userInfo.given_name} {userInfo.family_name[0]}.
            </Text>
          )}
        </Flex>
      </Button>
    ),
    [userInfo]
  );

  return (
    <Flex width={1} borderBottom="1px solid" borderColor="grey100" py={8}>
      <Breadcrumbs />
      <IconButton variant="default" mr={6} ml="auto" flex="0 0 auto" arial-label="Notifications">
        <Icon size="small" type="notification" />
      </IconButton>
      <Dropdown trigger={userButton} minWidth="100%">
        <Dropdown.Item onSelect={() => showSidesheet({ sidesheet: SIDESHEETS.EDIT_ACCOUNT })}>
          <MenuItem variant="default">Edit Profile</MenuItem>
        </Dropdown.Item>
        <Dropdown.Item onSelect={() => signOut({ onError: alert })}>
          <MenuItem variant="default" m={0}>
            Logout
          </MenuItem>
        </Dropdown.Item>
      </Dropdown>
    </Flex>
  );
};

export default Header;
