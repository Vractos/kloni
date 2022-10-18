import React from 'react'
import { Route, Routes } from "react-router-dom";
import AuthGuard from '../components/guards/authGuard';
import { routePaths } from '../constants/routes';
import Home from './home';
import Login from './login';

const Pages = () => {
  return (
    <Routes>
      <Route path={routePaths.LOGIN} element={<Login/>}/>
      <Route path={routePaths.HOME} element={<AuthGuard component={Home}/>}/>
    </Routes>
  )
}

export default Pages