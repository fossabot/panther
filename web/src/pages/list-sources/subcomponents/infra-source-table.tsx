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
