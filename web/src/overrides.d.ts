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
