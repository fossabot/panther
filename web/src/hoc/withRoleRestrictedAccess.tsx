import React from 'react';
import RoleRestrictedAccess, { RoleRestrictedAccessProps } from 'Components/role-restricted-access';

function withRoleRestrictedAccess<P>(config: Omit<RoleRestrictedAccessProps, 'children'>) {
  return (Component: React.FC<P>) => {
    const RoleRestrictedComponent: React.FC<P> = props => (
      <RoleRestrictedAccess {...config}>
        <Component {...props} />
      </RoleRestrictedAccess>
    );

    return RoleRestrictedComponent;
  };
}

export default withRoleRestrictedAccess;
