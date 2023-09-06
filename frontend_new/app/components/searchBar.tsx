import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';

export default function SearchBar() {
  return (
    <>
      <form>
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
            id="search"
            className="outline-none block p-4 pl-10 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500"
            placeholder="SKU"
            required
          />
          <button
            type="submit"
            className="text-white absolute right-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2"
          >
            Buscar
          </button>
        </div>
      </form>
    </>
  );
}
