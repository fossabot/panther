import { ErrorResponse } from 'apollo-link-error';
import storage from 'Helpers/storage';
import { UserInfo } from 'Components/utils/auth-context';
import { USER_INFO_STORAGE_KEY } from 'Source/constants';
import { Operation } from '@apollo/client';

interface ErrorData {
  operation?: Operation;
  extras?: {
    [key: string]: any;
  };
}

/**
 * Logs an error to sentry. Accepts *optional* additional arguments for easier debugging
 */
export const logError = (error: Error | ErrorResponse, { operation, extras }: ErrorData = {}) => {
  // On some environments we have sentry disabled
  const sentryDsn = process.env.SENTRY_DSN;
  const sentryRelease = process.env.PANTHER_VERSION;
  if (!sentryDsn) {
    return;
  }

  import(/* webpackChunkName: "sentry" */ '@sentry/browser').then(Sentry => {
    // We don't wanna initialize before any error occurs so we don't have to un-necessarily download
    // the sentry chunk at the user's device. `Init` method is idempotent, meaning that no matter
    // how many times we call it, it won't override anything. In addition it adds 0 thread overhead.
    Sentry.init({ dsn: sentryDsn, release: sentryRelease });
    // As soon as sentry is init, we add a scope to the error. Adding the scope here makes sure that
    // we don't have to manage the scopes on login/logout events
    Sentry.withScope(scope => {
      // Set the organization data and the email of the user
      const storedUserInfo = storage.read<UserInfo>(USER_INFO_STORAGE_KEY); // prettier-ignore
      if (storedUserInfo) {
        scope.setUser(storedUserInfo);
      }

      // If we have access to the operation that occurred, then we store this info for easier debugging
      if (operation) {
        scope.setTag('operationName', operation.operationName);
        scope.setExtra('operationVariables', operation.variables);
      }

      // If we have a custom stacktrace to share we add it here
      if (extras) {
        scope.setExtras(extras);
      }

      // Log the actual error
      Sentry.captureException(error);
    });
  });
};
