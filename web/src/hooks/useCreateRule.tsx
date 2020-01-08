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
import { MutationTuple } from '@apollo/client';
import useRouter from 'Hooks/useRouter';

interface UseCreateRuleProps<T> {
  mutation: MutationTuple<T, { [key: string]: any }>;
  getRedirectUri: (data: T) => string;
}

function useCreateRule<T>({ mutation, getRedirectUri }: UseCreateRuleProps<T>) {
  const { history } = useRouter();
  const [createRule, { data, error }] = mutation;

  const handleSubmit = React.useCallback(async values => {
    await createRule({ variables: { input: values } });
  }, []);

  React.useEffect(() => {
    if (data) {
      // After all is ok, navigate to the newly created resource
      history.push(getRedirectUri(data));
    }
  }, [data]);

  return { handleSubmit, data, error };
}

export default useCreateRule;
