/* A hook for getting access to the context value */
import React from 'react';
import { SidesheetContext } from 'Components/utils/sidesheet-context';

const useSidesheet = () => React.useContext(SidesheetContext);

export default useSidesheet;
