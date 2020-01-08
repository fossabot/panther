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
import { useQuery, gql } from '@apollo/client';
import { Integration } from 'Generated/schema';
import columns from 'Pages/list-sources/infra-source-columns';
import { INTEGRATION_TYPES } from 'Source/constants';
import BaseSourceTable from 'Pages/list-sources/subcomponents/base-source-table';

export const LIST_INFRA_SOURCES = gql`
  query ListInfraSources {
    integrations(input: { integrationType: "${INTEGRATION_TYPES.AWS_INFRA}" }) {
        awsAccountId
        createdAtTime
        createdBy
        integrationId
        integrationLabel
        integrationType
        scanEnabled
        scanIntervalMins
        scanStatus
        lastScanEndTime
      }
  }
`;

const InfraSourceTable = () => {
  const query = useQuery<{ integrations: Integration[] }>(LIST_INFRA_SOURCES, {
    fetchPolicy: 'cache-and-network',
  });

  return <BaseSourceTable query={query} columns={columns} />;
};

export default React.memo(InfraSourceTable);
