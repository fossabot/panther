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

/* The component responsible for rendering the actual sidesheets */
import React from 'react';
import useSidesheet from 'Hooks/useSidesheet';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import UpdateAwsSourcesSidesheet from 'Components/sidesheets/update-source-sidesheet';
import PolicyBulkUploadSidesheet from 'Components/sidesheets/policy-bulk-upload-sidesheet';
import SelectDestinationSidesheet from 'Components/sidesheets/select-destination-sidesheet';
import AddDestinationSidesheet from 'Components/sidesheets/add-destination-sidesheet';
import UpdateDestinationSidesheet from 'Components/sidesheets/update-destination-sidesheet';
import EditAccountSidesheet from 'Components/sidesheets/edit-account-sidesheet';
import UserInvitationSidesheet from 'Components/sidesheets/user-invitation-sidesheet';

const SidesheetManager: React.FC = () => {
  const { state: sidesheetState } = useSidesheet();
  if (!sidesheetState.sidesheet) {
    return null;
  }

  let Component;
  switch (sidesheetState.sidesheet) {
    case SIDESHEETS.ADD_DESTINATION:
      Component = AddDestinationSidesheet;
      break;
    case SIDESHEETS.UPDATE_DESTINATION:
      Component = UpdateDestinationSidesheet;
      break;
    case SIDESHEETS.SELECT_DESTINATION:
      Component = SelectDestinationSidesheet;
      break;
    case SIDESHEETS.UPDATE_SOURCE:
      Component = UpdateAwsSourcesSidesheet;
      break;
    case SIDESHEETS.POLICY_BULK_UPLOAD:
      Component = PolicyBulkUploadSidesheet;
      break;
    case SIDESHEETS.EDIT_ACCOUNT:
      Component = EditAccountSidesheet;
      break;
    case SIDESHEETS.USER_INVITATION:
      Component = UserInvitationSidesheet;
      break;
    default:
      break;
  }

  return <Component {...sidesheetState.props} />;
};

export default SidesheetManager;
