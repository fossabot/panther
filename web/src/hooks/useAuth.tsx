import React from 'react';
import { AuthContext } from 'Components/utils/auth-context';

const useAuth = () => React.useContext(AuthContext);

export default useAuth;
