
export default async function TableSkeleton() {
  return (
    <div className="overflow-x-auto relative rounded-lg">
      <table className="w-full text-sm text-left text-gray-500 table-auto">
        <thead className="text-xs text-gray-700 uppercase bg-gray-50">
          <tr>
            <th scope="col" className="py-3 px-6 text-center">
              Imagem
            </th>
            <th scope="col" className="py-3 px-6">
              Título do anúncio
            </th>
            <th scope="col" className="py-3 px-6">
              SKU
            </th>
            <th scope="col" className="py-3 px-6 text-center">
              Quantidade
            </th>
            <th scope="col" className="py-3 px-6 text-center">
              Preço
            </th>
            <th scope="col" className="py-3 px-6 text-center"></th>
          </tr>
        </thead>
        <tbody>
          {Array(5)
            .fill(0)
            .map((_, index) => {
              return (
                <tr className="bg-white border-b" key={index}>
                  <th
                    scope="row"
                    className="py-4 px-6 font-medium text-gray-900 whitespace-nowrap"
                  >
                    <div className="flex items-center justify-center">
                      <div className="w-10 h-10 flex-shrink-0">
                        <div className="rounded-full w-10 h-10 animate-pulse bg-gray-300" />
                      </div>
                    </div>
                  </th>
                  <td className="py-4 px-6">
                    <div className="h-2.5 rounded-full animate-pulse  bg-gray-300 w-[370px]" />
                  </td>
                  <td className="py-4 px-6">
                    <div className="h-2.5 rounded-full animate-pulse bg-gray-300 w-[119px]" />
                  </td>
                  <td className="py-4 px-6 text-center">
                    <div className="h-2.5 rounded-full animate-pulse bg-gray-300 w-5 ml-[29px]" />
                  </td>
                  <td className="py-4 px-6 text-center">
                    <div className="h-2.5 rounded-full animate-pulse  bg-gray-300 w-24" />
                  </td>
                  <td className="py-4 px-6 text-center">
                    <button
                      type="button"
                      className="group relative flex w-10/12 justify-center rounded-md border border-transparent animate-pulse bg-gray-300 py-2 px-4 text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                    >
                      <span className="absolute inset-y-0 left-0 flex items-center pl-3"></span>
                      <div className="h-4 rounded-full animate-pulse  bg-gray-300 w-11" />
                    </button>
                  </td>
                </tr>
              );
            })}
        </tbody>
      </table>
    </div>
  );
}
