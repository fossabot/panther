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
