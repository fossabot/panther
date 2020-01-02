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
