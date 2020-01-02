import React from 'react';
import { Box, Grid, Table, Alert } from 'pouncejs';
import { useQuery, gql } from '@apollo/client';
import Panel from 'Components/panel';
import urls from 'Source/urls';
import useRouter from 'Hooks/useRouter';
import {
  PolicySummary,
  ResourceSummary,
  OrganizationStatsResponse,
  IntegrationsByOrganizationResponse,
} from 'Generated/schema';
import ErrorBoundary from 'Components/error-boundary';
import { extractErrorMessage } from 'Helpers/utils';
import { topFailingPoliciesColumns, topFailingResourcesColumns } from './columns';
import PoliciesBySeverityChart from './subcomponents/policies-by-severity-chart';
import PoliciesByStatusChart from './subcomponents/policies-by-status-chart';
import ResourcesByPlatformChart from './subcomponents/resources-by-platform-chart';
import ResourcesByStatusChart from './subcomponents/resources-by-status-chart';
import DonutChartWrapper from './subcomponents/donut-chart-wrapper';
import OverviewPageEmptyDataFallback from './empty-data-fallback';
import OverviewPageSkeleton from './skeleton';

const GET_ORGANIZATION_STATS = gql`
  query GetOrganizationStats {
    organizationStats {
      scannedResources {
        byType {
          type
          count {
            fail
            pass
            error
          }
        }
      }
      appliedPolicies {
        info {
          error
          pass
          fail
        }
        low {
          error
          pass
          fail
        }
        medium {
          error
          pass
          fail
        }
        high {
          error
          pass
          fail
        }
        critical {
          error
          pass
          fail
        }
      }
      topFailingPolicies {
        id
        severity
      }
      topFailingResources {
        id
      }
    }
    integrations {
      integrations {
        integrationId
      }
    }
  }
`;

export type TopFailingPolicy = Pick<PolicySummary, 'id' | 'severity'>;
export type TopFailingResource = Pick<ResourceSummary, 'id'>;

interface ApolloQueryData {
  organizationStats: OrganizationStatsResponse;
  integrations: IntegrationsByOrganizationResponse;
}

const Overview: React.FC = () => {
  const { history } = useRouter();
  const { data, loading, error } = useQuery<ApolloQueryData>(GET_ORGANIZATION_STATS, {
    fetchPolicy: 'cache-and-network',
  });

  if (loading && !data) {
    return <OverviewPageSkeleton />;
  }

  if (error) {
    return (
      <Alert
        variant="error"
        title="We can\'t display this content right now"
        description={extractErrorMessage(error)}
      />
    );
  }

  if (!data.integrations.integrations.length) {
    return <OverviewPageEmptyDataFallback />;
  }

  return (
    <Box is="article" mb={6}>
      <Alert
        variant="info"
        title="Only active data is shown"
        description="Charts only include enabled policies which scanned at least one (1) resource & policies with a least one (1) policy attached to them"
        discardable
        mb={6}
      />
      <Grid
        gridTemplateColumns="repeat(4, 1fr)"
        gridRowGap={3}
        gridColumnGap={3}
        is="section"
        mb={3}
      >
        <DonutChartWrapper title="Policy Overview" icon="policy">
          <PoliciesBySeverityChart policies={data.organizationStats.appliedPolicies} />
        </DonutChartWrapper>
        <DonutChartWrapper title="Policy Failure Breakdown" icon="policy">
          <PoliciesByStatusChart policies={data.organizationStats.appliedPolicies} />
        </DonutChartWrapper>
        <DonutChartWrapper title="Resources Platforms" icon="resource">
          <ResourcesByPlatformChart resources={data.organizationStats.scannedResources} />
        </DonutChartWrapper>
        <DonutChartWrapper title="Resources Health" icon="resource">
          <ResourcesByStatusChart resources={data.organizationStats.scannedResources} />
        </DonutChartWrapper>
      </Grid>
      <Grid gridTemplateColumns="1fr 1fr" gridRowGap={2} gridColumnGap={3}>
        <Panel title="Top Failing Policies" size="small">
          <Box m={-6}>
            <ErrorBoundary>
              <Table<TopFailingPolicy>
                columns={topFailingPoliciesColumns}
                items={data.organizationStats.topFailingPolicies}
                getItemKey={policy => policy.id}
                onSelect={policy => history.push(urls.policies.details(policy.id))}
              />
            </ErrorBoundary>
          </Box>
        </Panel>
        <Panel title="Top Failing Resources" size="small">
          <Box m={-6}>
            <ErrorBoundary>
              <Table<TopFailingResource>
                columns={topFailingResourcesColumns}
                items={data.organizationStats.topFailingResources}
                getItemKey={resource => resource.id}
                onSelect={resource => history.push(urls.resources.details(resource.id))}
              />
            </ErrorBoundary>
          </Box>
        </Panel>
      </Grid>
    </Box>
  );
};

export default Overview;
