import { ApolloClient, ApolloLink, createHttpLink, InMemoryCache } from '@apollo/client';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import { History } from 'history';
import { createAuthLink, AUTH_TYPE } from 'aws-appsync-auth-link';
import { ErrorResponse, onError } from 'apollo-link-error';
import Auth from '@aws-amplify/auth';
import { LocationErrorState } from 'Components/utils/api-error-fallback';
import { LIST_REMEDIATIONS } from 'Components/forms/policy-form/policy-form-auto-remediation-fields';
import { logError } from 'Helpers/loggers';

/**
 * A link to react to GraphQL and/or network errors
 */
const createErrorLink = (history: History<LocationErrorState>) => {
  // Define the operations that won't trigger any handler actions or be logged anywhere (those can
  // still be handled by the component independently)
  const silentFailingOperations = [getOperationName(LIST_REMEDIATIONS)];

  return (onError(({ graphQLErrors, networkError, operation }: ErrorResponse) => {
    // If the error is not considered a fail, then don't log it to sentry
    if (silentFailingOperations.includes(operation.operationName)) {
      return;
    }

    if (graphQLErrors) {
      graphQLErrors.forEach(error => {
        logError(error, { operation });
        history.replace(history.location.pathname, { errorType: error.errorType });
      });
    }

    if (networkError) {
      logError(networkError, { operation });
    }
  }) as unknown) as ApolloLink;
};

/**
 * Typical HTTP link to add the GraphQL URL to query
 */
const httpLink = createHttpLink({ uri: process.env.GRAPHQL_ENDPOINT });

/**
 * This link is here to add the necessary headers present for AMAZON_COGNITO_USER_POOLS
 * authentication. It essentially signs the Authorization header with a JWT token
 */
const authLink = (createAuthLink({
  region: process.env.AWS_REGION,
  url: process.env.GRAPHQL_ENDPOINT,
  auth: {
    jwtToken: () => Auth.currentSession().then(session => session.getIdToken().getJwtToken()),
    type: AUTH_TYPE.AMAZON_COGNITO_USER_POOLS,
  },
}) as unknown) as ApolloLink;

/**
 * A function that will create an ApolloClient given a specific instance of a history
 */
const createApolloClient = (history: History<LocationErrorState>) =>
  new ApolloClient({
    link: ApolloLink.from([createErrorLink(history), authLink, httpLink]),
    cache: new InMemoryCache({
      typePolicies: {
        Destination: {
          keyFields: ['outputId'],
        },
      },
    }),
  });

export default createApolloClient;
