import React from 'react';
import useRouter from 'Hooks/useRouter';
import queryString from 'query-string';
import omitBy from 'lodash-es/omitBy';

const queryStringOptions = {
  arrayFormat: 'bracket' as const,
  parseNumbers: true,
  parseBooleans: true,
};

function useUrlParams<T extends { [key: string]: any }>() {
  const { history } = useRouter();

  /**
   * parses the query params of a URL and returns an object with params in the correct typo
   */
  const urlParams = queryString.parse(history.location.search, queryStringOptions) as T;

  /**
   * stringifies an object and adds it to the existing query params of a URL
   */
  const updateUrlParams = (params: Partial<T>) => {
    const mergedQueryParams = {
      ...urlParams,
      ...params,
    };

    // Remove any falsy value apart from the value `0` (number) and the value `false` (boolean)
    const cleanedMergedQueryParams = omitBy(
      mergedQueryParams,
      v => !v && !['number', 'boolean'].includes(typeof v)
    );

    history.replace(
      `${history.location.pathname}?${queryString.stringify(
        cleanedMergedQueryParams,
        queryStringOptions
      )}`
    );
  };

  // Cache those values as long as URL parameters are the same
  return React.useMemo(
    () => ({
      urlParams,
      updateUrlParams,
    }),
    [history.location.search]
  );
}

export default useUrlParams;
