import { useAuth0, withAuthenticationRequired } from '@auth0/auth0-react'
import React, { useEffect } from 'react'
import Loading from '../loading';

interface IGuardProps {
  component:  React.ComponentType<object>
}

const AuthGuard:React.FC<IGuardProps> = ({ component }) => {
  const Component = withAuthenticationRequired(component, {
    onRedirecting: () => <Loading/>,
  });

  return <Component />;
}

export default AuthGuard