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
import { SnackbarProvider, ThemeProvider } from 'pouncejs';
import { Router } from 'react-router-dom';
import Routes from 'Source/routes';
import { History } from 'history';
import { ApolloProvider } from '@apollo/client';
import { AuthProvider } from 'Components/utils/auth-context';
import { ModalProvider } from 'Components/utils/modal-context';
import { SidesheetProvider } from 'Components/utils/sidesheet-context';
import ModalManager from 'Components/utils/modal-manager';
import SidesheetManager from 'Components/utils/sidesheet-manager';
import ErrorBoundary from 'Components/error-boundary';
import createApolloClient from 'Source/client';

interface AppProps {
  history: History;
}
const App: React.FC<AppProps> = ({ history }) => {
  const client = React.useMemo(() => createApolloClient(history), [history]);
  return (
    <ErrorBoundary fallbackStrategy="passthrough">
      <ApolloProvider client={client}>
        <AuthProvider>
          <ThemeProvider>
            <Router history={history}>
              <SidesheetProvider>
                <ModalProvider>
                  <SnackbarProvider>
                    <Routes />
                    <ModalManager />
                    <SidesheetManager />
                  </SnackbarProvider>
                </ModalProvider>
              </SidesheetProvider>
            </Router>
          </ThemeProvider>
        </AuthProvider>
      </ApolloProvider>
    </ErrorBoundary>
  );
};

export default App;
