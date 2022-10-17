import { useAuth0 } from "@auth0/auth0-react";
import React from "react";
import Logo from '../../components/logo';
import { ArrowRightOnRectangleIcon } from "@heroicons/react/20/solid"

const Login: React.FC = () => {
  const { loginWithRedirect } = useAuth0();

  return (
    <main className='flex min-h-full items-center justify-center py-12 px-4 sm:px-6 lg:px-8'>
      <div className='w-full max-w-md space-y-8'>
        <div>
          <Logo className='mx-auto h-72 w-auto' />
        </div>
        <div>
          <button
            type="button"
            className="group relative flex w-full justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
            onClick={() => loginWithRedirect()}
          >
            <span className="absolute inset-y-0 left-0 flex items-center pl-3">
              <ArrowRightOnRectangleIcon className='h-5 w-5 text-indigo-500 group-hover:text-indigo-400' />
            </span>
            Entrar
          </button>
        </div>
      </div>
    </main>
  );

};

export default Login;