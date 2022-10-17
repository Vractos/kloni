import React from "react";
import sheepLogo from "../../assets/img/sheep.png"
export interface LogoProps {
  img?: string;
  className?: string;
}

const Logo: React.FC<LogoProps> = ({img = sheepLogo, className}) => {
  return <img src={img} alt="logo" className={className}/>
};

export default Logo;