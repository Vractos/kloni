import { withAuthenticationRequired } from '@auth0/auth0-react'
import React from 'react'

interface IGuardProps {
  component:  React.ComponentType<object>
}

const AuthGuard:React.FC<IGuardProps> = ({ component }) => {
  const Component = withAuthenticationRequired(component);

  return <Component />;
}

export default AuthGuard