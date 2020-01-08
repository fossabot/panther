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
import columns from 'Pages/list-sources/log-source-columns';
import { INTEGRATION_TYPES } from 'Source/constants';
import BaseSourceTable from 'Pages/list-sources/subcomponents/base-source-table';

export const LIST_LOG_SOURCES = gql`
  query ListLogSources {
    integrations(input: { integrationType: "${INTEGRATION_TYPES.AWS_LOGS}" }) {
          awsAccountId
          createdAtTime
          integrationId
          integrationLabel
          integrationType
          sourceSnsTopicArn
          logProcessingRoleArn
    }
  }
`;

const LogSourceTable = () => {
  const query = useQuery<{ integrations: Integration[] }>(LIST_LOG_SOURCES, {
    fetchPolicy: 'cache-and-network',
  });

  return <BaseSourceTable query={query} columns={columns} />;
};

export default React.memo(LogSourceTable);
