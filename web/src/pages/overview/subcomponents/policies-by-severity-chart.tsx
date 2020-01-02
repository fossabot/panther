import React from 'react';
import { defaultTheme } from 'pouncejs';
import { capitalize, countPoliciesBySeverityAndStatus } from 'Helpers/utils';
import DonutChart from 'Components/donut-chart';
import map from 'lodash-es/map';
import { OrganizationReportBySeverity } from 'Generated/schema';

const severityToGrayscaleMapping: {
  [key in keyof OrganizationReportBySeverity]: keyof typeof defaultTheme['colors'];
} = {
  critical: 'grey500',
  high: 'grey400',
  medium: 'grey300',
  low: 'grey200',
  info: 'grey100',
};

interface PoliciesBySeverityChartData {
  policies: OrganizationReportBySeverity;
}

const PoliciesBySeverityChart: React.FC<PoliciesBySeverityChartData> = ({ policies }) => {
  const allPoliciesChartData = map(
    severityToGrayscaleMapping,
    (color, severity: keyof OrganizationReportBySeverity) => ({
      value: countPoliciesBySeverityAndStatus(policies, severity, ['fail', 'error', 'pass']),
      label: capitalize(severity),
      color,
    })
  );

  return (
    <DonutChart data={allPoliciesChartData} renderLabel={(data, index) => data[index].value} />
  );
};

export default React.memo(PoliciesBySeverityChart);
