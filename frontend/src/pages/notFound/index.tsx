import React, { useCallback } from 'react'
import { useNavigate } from 'react-router-dom';
import { routePaths } from '../../constants/routes';

const NotFound = () => {
  const navigate = useNavigate()

  const goHome = useCallback(() => {
    navigate(routePaths.HOME)
  }, [])

  return (
    <>

      <div className="flex min-h-full flex-col bg-white pt-16 pb-12">
        <main className="mx-auto flex w-full max-w-7xl flex-grow flex-col justify-center px-4 sm:px-6 lg:px-8">

          <div className="py-16">
            <div className="text-center">
              <p className="text-4xl font-bold tracking-tight text-indigo-600 sm:text-5xl">404</p>
              <h1 className="mt-2 text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">Página não encontrada</h1>
              <div className="mt-6">
                <button onClick={goHome} className="text-base font-medium text-indigo-600 hover:text-indigo-500">
                  Ir para home
                  <span aria-hidden="true"> &rarr;</span>
                </button>
              </div>
            </div>
          </div>
        </main>
      </div>
    </>
  )
}

export default NotFound;