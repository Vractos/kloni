import React from 'react'
import { Navigate, Route, Routes } from "react-router-dom";
import AuthGuard from '../components/guards/authGuard';
import { routePaths } from '../constants/routes';
import Home from './home';
import Login from './login';
import NotFound from './notFound';

const Pages = () => {
  return (
    <Routes>
      <Route path={routePaths.LOGIN} element={<Login/>}/>
      <Route path={routePaths.HOME} element={<AuthGuard component={Home}/>}/>
      <Route path={routePaths.HOME} element={<AuthGuard component={Home}/>}/>
      <Route path={routePaths.NOT_FOUND} element={<NotFound/>}/>
      <Route path="*" element={<Navigate to={routePaths.NOT_FOUND}/>}/>
    </Routes>
  )
}

export default Pages