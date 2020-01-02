import React from 'react';
import urls from 'Source/urls';
import { Button, Dropdown, Flex, Icon, MenuItem } from 'pouncejs';
import useRouter from 'Hooks/useRouter';
import useSidesheet from 'Hooks/useSidesheet';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import { READONLY_ROLES_ARRAY } from 'Source/constants';

const CreateButton: React.FC = () => {
  const { history } = useRouter();
  const { showSidesheet } = useSidesheet();

  return (
    <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
      <Dropdown
        width={1}
        trigger={
          <Button size="large" variant="primary" is="div">
            <Flex>
              <Icon type="add" size="small" mr={2} />
              Create new
            </Flex>
          </Button>
        }
      >
        <Dropdown.Item onSelect={() => history.push(urls.policies.create())}>
          <MenuItem variant="default">Single</MenuItem>
        </Dropdown.Item>
        <Dropdown.Item
          onSelect={() =>
            showSidesheet({
              sidesheet: SIDESHEETS.POLICY_BULK_UPLOAD,
              props: { type: 'policy' },
            })
          }
        >
          <MenuItem variant="default">Bulk Upload</MenuItem>
        </Dropdown.Item>
      </Dropdown>
    </RoleRestrictedAccess>
  );
};

export default CreateButton;
