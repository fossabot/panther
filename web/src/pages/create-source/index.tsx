import React from 'react';
import queryString from 'query-string';
import { ADMIN_ROLES_ARRAY, INTEGRATION_TYPES } from 'Source/constants';
import Page404 from 'Pages/404';
import useRouter from 'Hooks/useRouter';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import CreateInfraSource from './subcomponents/create-infra-source';
import CreateLogSource from './subcomponents/create-log-source';

const CreateSourcePage: React.FC = () => {
  const { location } = useRouter();
  const { type } = queryString.parse(location.search) as { type: INTEGRATION_TYPES };

  return (
    <RoleRestrictedAccess allowedRoles={ADMIN_ROLES_ARRAY} fallback={<Page404 />}>
      {type === INTEGRATION_TYPES.AWS_INFRA ? <CreateInfraSource /> : <CreateLogSource />}
    </RoleRestrictedAccess>
  );
};

export default CreateSourcePage;
