/**
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program [The enterprise software] is licensed under the terms of a commercial license
 * available from Panther Labs Inc ("Panther Commercial License") by contacting contact@runpanther.com.
 * All use, distribution, and/or modification of this software, whether commercial or non-commercial,
 * falls under the Panther Commercial License to the extent it is permitted.
 */

import React from 'react';
import { Dropdown, Icon, IconButton, MenuItem } from 'pouncejs';
import useRouter from 'Hooks/useRouter';
import { User } from 'Generated/schema';
import { ADMIN_ROLES_ARRAY } from 'Source/constants';
import useModal from 'Hooks/useModal';
import { MODALS } from 'Components/utils/modal-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';

interface ListUsersTableRowOptionsProps {
  user: User;
}

const ListUsersTableRowOptions: React.FC<ListUsersTableRowOptionsProps> = ({ user }) => {
  const { location, history } = useRouter();
  const { showModal } = useModal();

  return (
    <RoleRestrictedAccess allowedRoles={ADMIN_ROLES_ARRAY}>
      <Dropdown
        trigger={
          <IconButton is="div" variant="default" my={-2}>
            <Icon type="more" size="small" />
          </IconButton>
        }
      >
        <Dropdown.Item onSelect={() => history.push(`${location.pathname}/${user.id}/edit/`)}>
          <MenuItem variant="default">Edit</MenuItem>
        </Dropdown.Item>
        <Dropdown.Item
          onSelect={() =>
            showModal({
              modal: MODALS.DELETE_USER,
              props: { user },
            })
          }
        >
          <MenuItem variant="default">Delete</MenuItem>
        </Dropdown.Item>
      </Dropdown>
    </RoleRestrictedAccess>
  );
};

export default React.memo(ListUsersTableRowOptions);
