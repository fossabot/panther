/* A hook for getting access to the context value */
import React from 'react';
import { ModalContext } from 'Components/utils/modal-context';

const useModal = () => React.useContext(ModalContext);

export default useModal;
