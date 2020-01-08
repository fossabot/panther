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

/* The component responsible for rendering the actual modals */
import React from 'react';
import useModal from 'Hooks/useModal';
import { MODALS } from 'Components/utils/modal-context';
import DeletePolicyModal from 'Components/modals/delete-policy-modal';
import DeleteUserModal from 'Components/modals/delete-user-modal';
import DeleteSourceModal from 'Components/modals/delete-source-modal';
import DeleteDestinationModal from 'Components/modals/delete-destination-modal';
import DeleteRuleModal from 'Components/modals/delete-rule-modal';
import NetworkErrorModal from 'Components/modals/network-error-modal';

const ModalManager: React.FC = () => {
  const { state: modalState } = useModal();
  if (!modalState.modal) {
    return null;
  }
  let Component;
  switch (modalState.modal) {
    case MODALS.DELETE_SOURCE:
      Component = DeleteSourceModal;
      break;
    case MODALS.DELETE_USER:
      Component = DeleteUserModal;
      break;
    case MODALS.DELETE_RULE:
      Component = DeleteRuleModal;
      break;
    case MODALS.DELETE_DESTINATION:
      Component = DeleteDestinationModal;
      break;
    case MODALS.NETWORK_ERROR:
      Component = NetworkErrorModal;
      break;
    case MODALS.DELETE_POLICY:
    default:
      Component = DeletePolicyModal;
      break;
  }

  return <Component {...modalState.props} />;
};

export default ModalManager;
