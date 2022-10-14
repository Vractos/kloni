import { useAuth0 } from "@auth0/auth0-react";
import React from "react";
import logoUrl from '../../assets/img/sheep.png'


export interface LogoProps {
  img?: string;
  height?: string;
  width?: string;
}

const LoginButton: React.FC<LogoProps> = ({img = ""}) => {
  const { loginWithRedirect } = useAuth0();

  // return <button onClick={() => loginWithRedirect()}>Log In</button>;
  return (
    <>
      <img src={img} alt="logo" />
    </>
  );

};

export default LoginButton;