import React from 'react';
import { Dropdown, Icon, IconButton, MenuItem } from 'pouncejs';
import useSidesheet from 'Hooks/useSidesheet';
import { Integration } from 'Generated/schema';
import useModal from 'Hooks/useModal';
import { MODALS } from 'Components/utils/modal-context';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import { READONLY_ROLES_ARRAY } from 'Source/constants';

interface ListSourcesTableRowOptionsProps {
  source: Integration;
}

const ListSourcesTableRowOptions: React.FC<ListSourcesTableRowOptionsProps> = ({ source }) => {
  const { showModal } = useModal();
  const { showSidesheet } = useSidesheet();

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
      <Dropdown
        trigger={
          <IconButton is="div" variant="default" my={-2}>
            <Icon type="more" size="small" />
          </IconButton>
        }
      >
        <Dropdown.Item
          onSelect={() =>
            showSidesheet({
              sidesheet: SIDESHEETS.UPDATE_SOURCE,
              props: { source },
            })
          }
        >
          <MenuItem variant="default">Edit</MenuItem>
        </Dropdown.Item>
        <Dropdown.Item
          onSelect={() =>
            showModal({
              modal: MODALS.DELETE_SOURCE,
              props: { source },
            })
          }
        >
          <MenuItem variant="default">Delete</MenuItem>
        </Dropdown.Item>
      </Dropdown>
    </RoleRestrictedAccess>
  );
};

export default React.memo(ListSourcesTableRowOptions);
