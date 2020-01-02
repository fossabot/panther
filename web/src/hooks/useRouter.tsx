import React from 'react';
import { __RouterContext as RouterContext, RouteComponentProps } from 'react-router';

function useRouter<Params = {}, StateShape = {}>() {
  return React.useContext(RouterContext) as RouteComponentProps<Params, undefined, StateShape>;
}

export default useRouter;
