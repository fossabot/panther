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
