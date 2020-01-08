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

import { useState, Dispatch, SetStateAction } from 'react';
import { useInfiniteScroll } from 'react-infinite-scroll-hook';

interface UseInfiniteScrollHookProps {
  loading: boolean;

  // The callback function to execute when we want to load more data.
  onLoadMore: Function;
}

// This hook builds upon https://www.npmjs.com/package/react-infinite-scroll-hook
const useInfiniteScrollHook = ({
  loading,
  onLoadMore,
}: UseInfiniteScrollHookProps): [
  React.MutableRefObject<undefined>,
  Dispatch<SetStateAction<boolean>>
] => {
  const [hasNextPage, setHasNextPage] = useState(true);
  const infiniteRef = useInfiniteScroll({
    loading,
    hasNextPage,
    onLoadMore,
    scrollContainer: 'window', // Set the scroll container to 'window' since 'parent' is a bit buggy
    checkInterval: 800, // The default is 200 which seems a bit too quick
  });

  return [infiniteRef, setHasNextPage];
};

export default useInfiniteScrollHook;
