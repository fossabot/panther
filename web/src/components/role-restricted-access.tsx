import React from 'react';
import { RoleNameEnum } from 'Generated/schema';
import useAuth from 'Hooks/useAuth';

export interface RoleRestrictedAccessProps {
  allowedRoles?: RoleNameEnum[];
  deniedRoles?: RoleNameEnum[];
  fallback?: React.ReactElement | null;
  children: React.ReactNode; // we need to specify it due to React.memo(..) down below
}

const RoleRestrictedAccess: React.FC<RoleRestrictedAccessProps> = ({
  allowedRoles,
  deniedRoles,
  fallback = null,
  children,
}) => {
  const { userInfo } = useAuth();

  if (!allowedRoles && !deniedRoles) {
    throw new Error(
      'You should specify either some roles to access the content or some to deny access to'
    );
  }

  if (!userInfo) {
    return fallback;
  }

  if (allowedRoles && userInfo.roles.some(role => allowedRoles.includes(role))) {
    return children as React.ReactElement;
  }

  if (deniedRoles && !userInfo.roles.every(role => deniedRoles.includes(role))) {
    return children as React.ReactElement;
  }

  return fallback;
};

export default React.memo(RoleRestrictedAccess);
