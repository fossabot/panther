import React from 'react';
import { DeletePolicyModalProps } from 'Components/modals/delete-policy-modal';
import { DeleteUserModalProps } from 'Components/modals/delete-user-modal';
import { DeleteSourceModalProps } from 'Components/modals/delete-source-modal';
import { DeleteDestinationModalProps } from 'Components/modals/delete-destination-modal';
import { DeleteRuleModalProps } from 'Components/modals/delete-rule-modal';

const SHOW_MODAL = 'SHOW_MODAL';
const HIDE_MODAL = 'HIDE_MODAL';

/* The available list of modals to dispatch */
export enum MODALS {
  DELETE_POLICY = 'DELETE_POLICY',
  DELETE_RULE = 'DELETE_RULE',
  DELETE_USER = 'DELETE_USER',
  DELETE_SOURCE = 'DELETE_SOURCE',
  DELETE_DESTINATION = 'DELETE_DESTINATION',
  NETWORK_ERROR = 'NETWORK_ERROR',
}

/* The shape of the reducer state */
interface ModalStateShape {
  modal: keyof typeof MODALS | null;
  props: { [key: string]: any };
}

/* 1st action */
interface ShowPolicyModalAction {
  type: typeof SHOW_MODAL;
  payload: {
    modal: MODALS.DELETE_POLICY;
    props: DeletePolicyModalProps;
  };
}

/* 2nd action */
interface HideModalAction {
  type: typeof HIDE_MODAL;
}

/* Delete User action */
interface ShowDeleteUserModalAction {
  type: typeof SHOW_MODAL;
  payload: {
    modal: MODALS.DELETE_USER;
    props: DeleteUserModalProps;
  };
}

/* Delete Source action */
interface ShowDeleteSourceModalAction {
  type: typeof SHOW_MODAL;
  payload: {
    modal: MODALS.DELETE_SOURCE;
    props: DeleteSourceModalProps;
  };
}

/* 1st action */
interface ShowDeleteRuleModalAction {
  type: typeof SHOW_MODAL;
  payload: {
    modal: MODALS.DELETE_RULE;
    props: DeleteRuleModalProps;
  };
}

/* Delete Destination action */
interface ShowDeleteDestinationModalAction {
  type: typeof SHOW_MODAL;
  payload: {
    modal: MODALS.DELETE_DESTINATION;
    props: DeleteDestinationModalProps;
  };
}

/* Delete Destination action */
interface ShowNetworkErrorModalAction {
  type: typeof SHOW_MODAL;
  payload: {
    modal: MODALS.NETWORK_ERROR;
  };
}

/* The available actions that can be dispatched */
type ModalStateAction =
  | ShowDeleteSourceModalAction
  | ShowDeleteUserModalAction
  | ShowPolicyModalAction
  | ShowDeleteRuleModalAction
  | ShowDeleteDestinationModalAction
  | ShowNetworkErrorModalAction
  | HideModalAction;

/* initial state of the reducer */
const initialState: ModalStateShape = {
  modal: null,
  props: {},
};

const modalReducer = (state: ModalStateShape, action: ModalStateAction) => {
  switch (action.type) {
    case SHOW_MODAL:
      return {
        modal: action.payload.modal,
        props: 'props' in action.payload ? action.payload.props : {},
      };
    case HIDE_MODAL:
      return { modal: null, props: {} };
    default:
      return state;
  }
};

interface ModalContextValue {
  state: ModalStateShape;
  showModal: (input: Exclude<ModalStateAction, HideModalAction>['payload']) => void;
  hideModal: () => void;
}

/* Context that will hold the `state` and `dispatch` */
export const ModalContext = React.createContext<ModalContextValue>(undefined);

/* A enhanced version of the context provider */
export const ModalProvider: React.FC = ({ children }) => {
  const [state, dispatch] = React.useReducer<React.Reducer<ModalStateShape, ModalStateAction>>(
    modalReducer,
    initialState
  );

  // for perf reasons we only want to re-render on state updates
  const contextValue = React.useMemo(
    () => ({
      state,
      hideModal: () => dispatch({ type: 'HIDE_MODAL' }),
      showModal: ({ modal, props }) => dispatch({ type: 'SHOW_MODAL', payload: { modal, props } }),
    }),
    [state]
  );

  // make the `state` and `dispatch` available to the components
  return <ModalContext.Provider value={contextValue}>{children}</ModalContext.Provider>;
};
