import 'yup';
import 'graphql';

/**
 * We declare a custom `unique` yup method and we want to expose it through the global package, so
 * that every module that imports yup can have access to it
 */
declare module 'yup' {
  export interface Schema<T> {
    unique<T>(message: string, key?: keyof T): this;
  }
}

/**
 * We are utilising AppSync, whose error doesn't conform to the standardized error set by GraphQL
 * itself (what a surprise). Thus, we need to add the fields that AppSync returns to the schema of
 * the GraphQL error
 */
declare module 'graphql' {
  export interface GraphQLError {
    errorType: string;
    errorInfo?: any;
  }
}
