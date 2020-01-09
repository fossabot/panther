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
import { capitalize, countPoliciesBySeverityAndStatus } from 'Helpers/utils';
import DonutChart from 'Components/donut-chart';
import map from 'lodash-es/map';
import sum from 'lodash-es/sum';
import { OrganizationReportBySeverity } from 'Generated/schema';
import { defaultTheme } from 'pouncejs';

const severityToColorMapping: {
  [key in keyof OrganizationReportBySeverity]: keyof typeof defaultTheme['colors'];
} = {
  critical: 'red300',
  high: 'red200',
  medium: 'blue100',
  low: 'grey100',
  info: 'grey50',
};

interface PoliciesByStatusChartData {
  policies: OrganizationReportBySeverity;
}

const PoliciesByStatusChart: React.FC<PoliciesByStatusChartData> = ({ policies }) => {
  const severities = Object.keys(severityToColorMapping);
  const totalPolicies = sum(
    severities.map((severity: keyof OrganizationReportBySeverity) =>
      countPoliciesBySeverityAndStatus(policies, severity, ['fail', 'error', 'pass'])
    )
  );

  const failingPoliciesChartData = [
    ...map(severityToColorMapping, (color, severity: keyof OrganizationReportBySeverity) => ({
      value: countPoliciesBySeverityAndStatus(policies, severity, ['fail', 'error']),
      label: capitalize(severity),
      color,
    })),
    {
      value: sum(
        Object.keys(severityToColorMapping).map((severity: keyof OrganizationReportBySeverity) =>
          countPoliciesBySeverityAndStatus(policies, severity, ['pass'])
        )
      ),
      label: 'Passing',
      color: 'green100' as const,
    },
  ];

  return (
    <DonutChart
      data={failingPoliciesChartData}
      renderLabel={(chartData, index) => {
        const { value: severityGroupingValue } = chartData[index];
        const percentage = Math.round((severityGroupingValue * 100) / totalPolicies).toFixed(0);

        return `${severityGroupingValue}\n{small|${percentage}% of all}`;
      }}
    />
  );
};

export default React.memo(PoliciesByStatusChart);