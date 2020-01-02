import React from 'react';
import useUrlParams from 'Hooks/useUrlParams';

function useRequestParamsWithPagination<AvailableParams extends { page?: number }>() {
  const { urlParams, updateUrlParams } = useUrlParams<Partial<AvailableParams>>();

  // This is our typical function that updates the parameters with the addition of resetting the
  // page to `1`. The only time where we don't wanna do that, is when changing pages. In this
  // scenario, we want to change the params but not reset the page.
  const updateRequestParamsAndResetPaging = React.useCallback(
    (newParams: Partial<AvailableParams>) => {
      updateUrlParams({ ...urlParams, ...newParams, page: 1 });
    },
    [urlParams]
  );

  // This is a similar function like the above but instead of updating the existing params with the
  // new parameters, it clears all the parameters and just sets the parameters passed as an argument
  const setRequestParamsAndResetPaging = React.useCallback(
    (newParams: Partial<AvailableParams>) => {
      updateUrlParams({ ...newParams, page: 1 });
    },
    [urlParams]
  );

  // This is the function to call whenever a page changes. The difference lies in the value of the
  // `page` value
  const updatePagingParams = React.useCallback(
    (newPage: AvailableParams['page']) => {
      updateUrlParams({ ...urlParams, page: newPage });
    },
    [urlParams]
  );

  return React.useMemo(
    () => ({
      requestParams: urlParams,
      updateRequestParamsAndResetPaging,
      setRequestParamsAndResetPaging,
      updatePagingParams,
    }),
    [urlParams]
  );
}
export default useRequestParamsWithPagination;
