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
import { UpdateAWSSourcesSidesheetProps } from 'Components/sidesheets/update-source-sidesheet';
import { AddDestinationSidesheetProps } from 'Components/sidesheets/add-destination-sidesheet';
import { UpdateDestinationSidesheetProps } from 'Components/sidesheets/update-destination-sidesheet';
import { PolicyBulkUploadSideSheetProps } from 'Components/sidesheets/policy-bulk-upload-sidesheet';

const SHOW_SIDESHEET = 'SHOW_SIDESHEET';
const HIDE_SIDESHEET = 'HIDE_SIDESHEET';

/* The available list of sidesheets to dispatch */
export enum SIDESHEETS {
  UPDATE_SOURCE = 'UPDATE_SOURCE',
  POLICY_BULK_UPLOAD = 'POLICY_BULK_UPLOAD',
  SELECT_DESTINATION = 'SELECT_DESTINATION',
  ADD_DESTINATION = 'ADD_DESTINATION',
  UPDATE_DESTINATION = 'UPDATE_DESTINATION',
  EDIT_ACCOUNT = 'EDIT_ACCOUNT',
  USER_INVITATION = 'USER_INVITATION',
}

/* The shape of the reducer state */
interface SidesheetStateShape {
  sidesheet: keyof typeof SIDESHEETS | null;
  props?: { [key: string]: any };
}

/* 1st action */
interface UpdateSourceSideSheetAction {
  type: typeof SHOW_SIDESHEET;
  payload: {
    sidesheet: SIDESHEETS.UPDATE_SOURCE;
    props: UpdateAWSSourcesSidesheetProps;
  };
}

interface SelectDestinationSideSheetAction {
  type: typeof SHOW_SIDESHEET;
  payload: {
    sidesheet: SIDESHEETS.SELECT_DESTINATION;
  };
}

interface AddDestinationSideSheetAction {
  type: typeof SHOW_SIDESHEET;
  payload: {
    sidesheet: SIDESHEETS.ADD_DESTINATION;
    props: AddDestinationSidesheetProps;
  };
}

interface UpdateDestinationSideSheetAction {
  type: typeof SHOW_SIDESHEET;
  payload: {
    sidesheet: SIDESHEETS.UPDATE_DESTINATION;
    props: UpdateDestinationSidesheetProps;
  };
}

interface HideSidesheetAction {
  type: typeof HIDE_SIDESHEET;
}

/* Bulk upload policies action */
interface PolicyBulkUploadSideSheetAction {
  type: typeof SHOW_SIDESHEET;
  payload: {
    sidesheet: SIDESHEETS.POLICY_BULK_UPLOAD;
    props: PolicyBulkUploadSideSheetProps;
  };
}

interface EditAccountSideSheetAction {
  type: typeof SHOW_SIDESHEET;
  payload: {
    sidesheet: SIDESHEETS.EDIT_ACCOUNT;
  };
}

interface UserInvitationSideSheetAction {
  type: typeof SHOW_SIDESHEET;
  payload: {
    sidesheet: SIDESHEETS.USER_INVITATION;
  };
}

/* The available actions that can be dispatched */
type SidesheetStateAction =
  | UpdateSourceSideSheetAction
  | PolicyBulkUploadSideSheetAction
  | SelectDestinationSideSheetAction
  | AddDestinationSideSheetAction
  | UpdateDestinationSideSheetAction
  | EditAccountSideSheetAction
  | UserInvitationSideSheetAction
  | HideSidesheetAction;

/* initial state of the reducer */
const initialState: SidesheetStateShape = {
  sidesheet: null,
  props: {},
};

const sidesheetReducer = (state: SidesheetStateShape, action: SidesheetStateAction) => {
  switch (action.type) {
    case SHOW_SIDESHEET:
      return {
        sidesheet: action.payload.sidesheet,
        props: 'props' in action.payload ? action.payload.props : {},
      };
    case HIDE_SIDESHEET:
      return { sidesheet: null, props: {} };
    default:
      return state;
  }
};

interface SidesheetContextValue {
  state: SidesheetStateShape;
  showSidesheet: (input: Exclude<SidesheetStateAction, HideSidesheetAction>['payload']) => void;
  hideSidesheet: () => void;
}

/* Context that will hold the `state` and `dispatch` */
export const SidesheetContext = React.createContext<SidesheetContextValue>(undefined);

/* A enhanced version of the context provider */
export const SidesheetProvider: React.FC = ({ children }) => {
  const [state, dispatch] = React.useReducer<
    React.Reducer<SidesheetStateShape, SidesheetStateAction>
  >(sidesheetReducer, initialState);

  // for perf reasons we only want to re-render on state updates
  const contextValue = React.useMemo(
    () => ({
      state,
      hideSidesheet: () => dispatch({ type: 'HIDE_SIDESHEET' }),
      showSidesheet: ({ sidesheet, props }) =>
        dispatch({ type: 'SHOW_SIDESHEET', payload: { sidesheet, props } }),
    }),
    [state]
  );

  // make the `state` and `dispatch` available to the components
  return <SidesheetContext.Provider value={contextValue}>{children}</SidesheetContext.Provider>;
};
