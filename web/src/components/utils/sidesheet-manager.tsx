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
