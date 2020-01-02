import React from 'react';
import { useQuery, gql } from '@apollo/client';
import { IntegrationsByOrganizationResponse } from 'Generated/schema';
import columns from 'Pages/list-sources/infra-source-columns';
import { INTEGRATION_TYPES } from 'Source/constants';
import BaseSourceTable from 'Pages/list-sources/subcomponents/base-source-table';

export const LIST_INFRA_SOURCES = gql`
  query ListInfraSources {
    integrations(input: { integrationType: "${INTEGRATION_TYPES.AWS_INFRA}" }) {
      integrations {
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
  }
`;

const InfraSourceTable = () => {
  const query = useQuery<{ integrations: IntegrationsByOrganizationResponse }>(LIST_INFRA_SOURCES, {
    fetchPolicy: 'cache-and-network',
  });

  return (
    <BaseSourceTable
      query={query}
      columns={columns}
      integrationType={INTEGRATION_TYPES.AWS_INFRA}
    />
  );
};

export default React.memo(InfraSourceTable);
