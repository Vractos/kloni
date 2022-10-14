import React from "react";

export interface LogoProps {
  img?: string;
  height?: string;
  width?: string;
}

const Logo: React.FC<LogoProps> = ({img = "../../assets/sheep.png"}) => {
  return <img src={img} alt="logo" />
};

export default Logo;