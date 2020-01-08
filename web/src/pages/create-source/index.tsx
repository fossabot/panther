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
