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

import { AlertSummary, PolicySummary, ResourceSummary, RuleSummary } from 'Generated/schema';
import { INTEGRATION_TYPES } from 'Source/constants';

const urls = {
  overview: () => '/overview',
  rules: {
    list: () => '/rules/',
    details: (id: RuleSummary['id']) => `${urls.rules.list()}${encodeURIComponent(id)}/`,
    edit: (id: RuleSummary['id']) => `${urls.rules.details(id)}edit/`,
    create: () => `${urls.rules.list()}new/`,
  },
  policies: {
    list: () => '/policies/',
    details: (id: PolicySummary['id']) => `${urls.policies.list()}${encodeURIComponent(id)}/`,
    edit: (id: PolicySummary['id']) => `${urls.policies.details(id)}edit/`,
    create: () => `${urls.policies.list()}new/`,
  },
  resources: {
    list: () => '/resources/',
    details: (id: ResourceSummary['id']) => `${urls.resources.list()}${encodeURIComponent(id)}/`,
    edit: (id: ResourceSummary['id']) => `${urls.resources.details(id)}edit/`,
  },
  alerts: {
    list: () => '/alerts/',
    details: (id: AlertSummary['alertId']) => `${urls.alerts.list()}${encodeURIComponent(id)}/`,
  },
  account: {
    settings: {
      overview: () => `/settings/`,
      general: () => `${urls.account.settings.overview()}general`,
      users: () => `${urls.account.settings.overview()}users`,
      sources: {
        list: () => `${urls.account.settings.overview()}sources/`,
        create: (integrationType?: INTEGRATION_TYPES) =>
          `${urls.account.settings.sources.list()}new/${
            integrationType ? `?type=${integrationType}` : ''
          }`,
      },
      destinations: () => `${urls.account.settings.overview()}destinations`,
    },

    auth: {
      signIn: () => `/sign-in/`,
      forgotPassword: () => `/password-forgot/`,
      resetPassword: () => `/password-reset/`,
    },
  },

  integrations: {
    details: (serviceName: string) => `/integrations/${serviceName}`,
  },
};

export default urls;
