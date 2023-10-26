export default function Integrations() {
  return (
    <>
      <main className="px-4 py-16 sm:px-6 lg:flex-auto lg:px-0 lg:py-20">
        <div className="mx-auto max-w-2xl space-y-16 sm:space-y-20 lg:mx-0 lg:max-w-none">
          <div>
            <h2 className="text-base font-semibold leading-7 text-gray-900">
              Integrações
            </h2>
            <p className="mt-1 text-sm leading-6 text-gray-500">
              Conecte suas contas do Mercado Livre.
            </p>

            <ul
              role="list"
              className="mt-6 divide-y divide-gray-100 border-t border-gray-200 text-sm leading-6"
            >
              <li className="flex justify-between gap-x-6 py-6">
                <div className="font-medium text-gray-900">Mercado Livre</div>
                <button
                  type="button"
                  className="font-semibold text-indigo-600 hover:text-indigo-500"
                >
                  Atualizar
                </button>
              </li>
            </ul>

            <div className="flex border-t border-gray-100 pt-6">
              <button
                type="button"
                disabled
                className="text-sm font-semibold leading-6 text-indigo-600 hover:text-indigo-500"
              >
                <span aria-hidden="true">+</span> Adicionar outra conta{" "}
                <span className="inline-flex items-center rounded-full bg-yellow-50 ml-1 px-1.5 py-0.5 text-xs font-medium text-yellow-800 ring-1 ring-inset ring-yellow-600/20">
                  Em breve
                </span>
              </button>
            </div>
            <div className="flex pt-12"></div>
          </div>
        </div>
      </main>
    </>
  );
}
