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
