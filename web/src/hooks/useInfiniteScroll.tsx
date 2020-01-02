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
