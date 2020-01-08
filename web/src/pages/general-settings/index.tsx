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
import { Alert, Box } from 'pouncejs';
import { useQuery, gql } from '@apollo/client';
import { ADMIN_ROLES_ARRAY } from 'Source/constants';
import { GetOrganizationResponse } from 'Generated/schema';
import CompanyInformation from 'Pages/general-settings/subcomponent/company-information-panel';
import RoleRestrictedAccess from 'Components/role-restricted-access';
import Page404 from 'Pages/404';
import ErrorBoundary from 'Components/error-boundary';
import { extractErrorMessage } from 'Helpers/utils';
import GeneralSettingsPageSkeleton from './skeleton';

export const GET_ORGANIZATION = gql`
  query GetOrganization {
    organization {
      organization {
        id
        displayName
        email
        alertReportFrequency
        remediationConfig {
          awsRemediationLambdaArn
        }
      }
    }
  }
`;

interface ApolloQueryData {
  organization: GetOrganizationResponse;
}

// Parent container for the general settings section
const GeneralSettingsContainer: React.FC = () => {
  // We're going to fetch the organization info at the top level and pass down relevant attributes and loading for each panel
  const { loading, error, data } = useQuery<ApolloQueryData>(GET_ORGANIZATION, {
    fetchPolicy: 'cache-and-network',
  });

  if (loading) {
    return <GeneralSettingsPageSkeleton />;
  }

  if (error) {
    return (
      <Alert
        variant="error"
        title="Failed to query company information"
        description={
          extractErrorMessage(error) ||
          'Sorry, something went wrong, please reach out to support@runpanther.io if this problem persists'
        }
      />
    );
  }

  return (
    <RoleRestrictedAccess allowedRoles={ADMIN_ROLES_ARRAY} fallback={<Page404 />}>
      <Box mb={6}>
        <ErrorBoundary>
          <CompanyInformation
            displayName={data.organization.organization.displayName}
            email={data.organization.organization.email}
          />
        </ErrorBoundary>
      </Box>
    </RoleRestrictedAccess>
  );
};

export default GeneralSettingsContainer;
