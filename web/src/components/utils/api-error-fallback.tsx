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
import useModal from 'Hooks/useModal';
import { MODALS } from 'Components/utils/modal-context';
import useRouter from 'Hooks/useRouter';
import Page404 from 'Pages/404';
import Page403 from 'Pages/403';
import Page500 from 'Pages/500';
import urls from 'Source/urls';
import { Redirect } from 'react-router-dom';

export interface LocationErrorState {
  errorType?: string;
}

const ApiErrorFallback: React.FC = ({ children }) => {
  const { location } = useRouter<{}, LocationErrorState>();
  const { showModal, hideModal } = useModal();

  const showNetworkErroModal = React.useCallback(() => {
    showModal({ modal: MODALS.NETWORK_ERROR });
  }, []);

  const hideNetworkErroModal = React.useCallback(() => {
    hideModal();
  }, []);

  React.useEffect(() => {
    window.addEventListener('offline', showNetworkErroModal);
    window.addEventListener('online', hideNetworkErroModal);

    return () => {
      window.removeEventListener('offline', showNetworkErroModal);
      window.removeEventListener('online', hideNetworkErroModal);
    };
  }, []);

  switch (location.state?.errorType) {
    case '401':
      return <Redirect to={urls.account.auth.signIn()} />;
    case '404':
      return <Page404 />;
    case '403':
      return <Page403 />;
    case '500':
      return <Page500 />;
    default:
      return children as React.ReactElement;
  }
};

export default ApiErrorFallback;
