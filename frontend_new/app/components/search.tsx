'use client';

import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import { useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState, } from 'react';
import { createUrl } from '../_lib/utils/url';


export default function Search() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [searchValue, setSearchValue] = useState('');

  useEffect(() => {
    setSearchValue(searchParams?.get('q') || '');
    console.log(searchParams)

  }, [searchParams, setSearchValue]);

  function onSubmit(e: React.FormEvent) {
    e.preventDefault();

    const val = e.target as HTMLFormElement;
    const search = val.search as HTMLInputElement;
    const newParams = new URLSearchParams(searchParams.toString());

    if (search.value) {
      newParams.set('q', search.value);
    } else {
      newParams.delete('q');
    }

    router.push(createUrl('/', newParams));
  }

  return (
      <form onSubmit={onSubmit}>
        <label
          htmlFor="search"
          className="mb-2 text-sm font-medium sr-only dark:text-gray-300"
        >
          Buscar
        </label>
        <div className="relative w-full ">
          <div className="flex absolute inset-y-0 left-0 items-center pl-3 pointer-events-none">
            <MagnifyingGlassIcon className="h-5 w-5 text-indigo-500 group-hover:text-indigo-400" />
          </div>
          <input
            type="search"
            name='search'
            placeholder="Utilize o SKU para buscar"
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            className="outline-none block p-4 pl-10 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-indigo-500 focus:border-indigo-500"
          />
          <button
            type="submit"
            onSubmit={e => onSubmit(e)}
            className="text-white absolute right-2.5 bottom-2.5 bg-indigo-600 hover:bg-indigo-700 focus:ring-4 focus:outline-none focus:ring-indigo-300 font-medium rounded-lg text-sm px-4 py-2"
          >
            Buscar
          </button>
        </div>
      </form>
  );
}
