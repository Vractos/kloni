"use client"; // Error components must be Client Components

import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";
import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error;
  reset: () => void;
}) {
  useEffect(() => {
    // Log the error to an error reporting service
    console.error(error);
  }, [error]);

  return (
    <>
      <div className="overflow-hidden rounded-lg bg-slate-200 mt-2">
        <div className="px-4 py-5 sm:p-6 text-center">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="mx-auto h-20 w-20 text-gray-400"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>

         
          <h3 className="mt-3 text-base font-semibold text-gray-900">
            Algo deu errado
          </h3>
          <p className="mt-3 text-xs text-gray-500">
            Não foi possível obter os anúncios.
          </p>
          <div className="mt-3">
            <button
              type="button"
              onClick={() => reset()}
              className="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
              Tentar novamente
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
