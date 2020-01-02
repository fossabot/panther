import React from 'react';
import { useQuery, gql } from '@apollo/client';
import { IntegrationsByOrganizationResponse } from 'Generated/schema';
import columns from 'Pages/list-sources/log-source-columns';
import { INTEGRATION_TYPES } from 'Source/constants';
import BaseSourceTable from 'Pages/list-sources/subcomponents/base-source-table';

export const LIST_LOG_SOURCES = gql`
  query ListLogSources {
    integrations(input: { integrationType: "${INTEGRATION_TYPES.AWS_LOGS}" }) {
      integrations {
          awsAccountId
          createdAtTime
          integrationId
          integrationLabel
          integrationType
          sourceSnsTopicArn
          logProcessingRoleArn
      }
    }
  }
`;

const LogSourceTable = () => {
  const query = useQuery<{ integrations: IntegrationsByOrganizationResponse }>(LIST_LOG_SOURCES, {
    fetchPolicy: 'cache-and-network',
  });

  return (
    <BaseSourceTable query={query} columns={columns} integrationType={INTEGRATION_TYPES.AWS_LOGS} />
  );
};

export default React.memo(LogSourceTable);
