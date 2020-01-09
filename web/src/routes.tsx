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
import { Redirect, Route, Switch } from 'react-router-dom';
import ListPoliciesPage from 'Pages/list-policies';
import OverviewPage from 'Pages/overview';
import ListResourcesPage from 'Pages/list-resources';
import ResourceDetailsPage from 'Pages/resource-details';
import PolicyDetailsPage from 'Pages/policy-details';
import GeneralSettingsPage from 'Pages/general-settings';
import SourcesPage from 'Pages/list-sources';
import CreateSourcesPage from 'Pages/create-source';
import SignInPage from 'Pages/sign-in';
import DestinationsPage from 'Pages/destinations';
import UsersPage from 'Pages/users';
import RuleDetailsPage from 'Pages/rule-details';
import ListRulesPage from 'Pages/list-rules';
import EditRulePage from 'Pages/edit-rule';
import CreateRulePage from 'Pages/create-rule';
import AlertDetailsPage from 'Pages/alert-details';
import EditPolicyPage from 'Pages/edit-policy';
import CreatePolicyPage from 'Pages/create-policy';
import ListAlertsPage from 'Pages/list-alerts';
import Layout from 'Components/layout';
import urls from 'Source/urls';
import GuardedRoute from 'Components/guarded-route';
import ForgotPasswordPage from 'Pages/forgot-password';
import ForgotPasswordConfirmPage from 'Pages/forgot-password-confirm';
import ErrorBoundary from 'Components/error-boundary';
import Page404 from 'Pages/404';
import APIErrorFallback from 'Components/utils/api-error-fallback';

// Main page container for the web application, Navigation bar and Content body goes here
const PrimaryPageLayout: React.FunctionComponent = () => {
  return (
    <Switch>
      <GuardedRoute
        limitAccessTo="anonymous"
        exact
        path={urls.account.auth.signIn()}
        component={SignInPage}
      />
      <GuardedRoute
        limitAccessTo="anonymous"
        exact
        path={urls.account.auth.forgotPassword()}
        component={ForgotPasswordPage}
      />
      <GuardedRoute
        limitAccessTo="anonymous"
        exact
        path={urls.account.auth.resetPassword()}
        component={ForgotPasswordConfirmPage}
      />
      <GuardedRoute path="/" limitAccessTo="authenticated">
        <Layout>
          <ErrorBoundary>
            <APIErrorFallback>
              <Switch>
                <Route exact path={urls.overview()} component={OverviewPage} />
                <Route exact path={urls.policies.list()} component={ListPoliciesPage} />
                <Route exact path={urls.policies.create()} component={CreatePolicyPage} />
                <Route exact path={`${urls.policies.list()}:id`} component={PolicyDetailsPage} />
                <Route exact path={`${urls.policies.list()}:id/edit`} component={EditPolicyPage} />
                <Route exact path={urls.resources.list()} component={ListResourcesPage} />
                <Route exact path={`${urls.resources.list()}:id`} component={ResourceDetailsPage} />
                <Route exact path={urls.account.settings.sources.list()} component={SourcesPage} />
                <Route
                  exact
                  path={urls.account.settings.general()}
                  component={GeneralSettingsPage}
                />
                <Route exact path={urls.account.settings.users()} component={UsersPage} />
                <Route
                  exact
                  path={urls.account.settings.destinations()}
                  component={DestinationsPage}
                />
                <Route
                  exact
                  path={urls.account.settings.sources.create()}
                  component={CreateSourcesPage}
                />
                <Route exact path={`${urls.alerts.list()}:id`} component={AlertDetailsPage} />
                <Route exact path={urls.rules.list()} component={ListRulesPage} />
                <Route exact path={urls.alerts.list()} component={ListAlertsPage} />
                <Route exact path={urls.rules.create()} component={CreateRulePage} />
                <Route exact path={`${urls.rules.list()}:id`} component={RuleDetailsPage} />
                <Route exact path={`${urls.rules.list()}:id/edit`} component={EditRulePage} />
                <Redirect exact from="/" to={urls.overview()} />
                <Route component={Page404} />
              </Switch>
            </APIErrorFallback>
          </ErrorBoundary>
        </Layout>
      </GuardedRoute>
    </Switch>
  );
};

export default PrimaryPageLayout;