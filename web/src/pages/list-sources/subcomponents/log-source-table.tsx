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
