import Image from "next/image";
import { formatCurrency } from "../../_lib/utils/formatter";
import { getAnnouncements } from "../../api/handlers/announcements";
import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

export default async function Table({ sku }: { sku: string }) {
  const announcements = await getAnnouncements(sku, "test");
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
          {announcements &&
            announcements.map((announcement, index) => {
              return (
                <tr className="bg-white border-b" key={index}>
                  <th
                    scope="row"
                    className="py-4 px-6 font-medium text-gray-900 whitespace-nowrap"
                  >
                    <div className="flex items-center justify-center">
                      <div className="w-10 h-10 flex-shrink-0">
                        <Image
                          className="rounded-full"
                          src={announcement.picture}
                          width={40}
                          height={40}
                          alt={announcement.title}
                        />
                      </div>
                    </div>
                  </th>
                  <td className="py-4 px-6">
                    <a
                      href={announcement.link}
                      target="_blank"
                      rel="noreferrer"
                      className="text-blue-600"
                    >
                      {announcement.title}
                    </a>
                  </td>
                  <td className="py-4 px-6">{announcement.sku}</td>
                  <td className="py-4 px-6 text-center">
                    {announcement.quantity}
                  </td>
                  <td className="py-4 px-6 text-center">
                    {formatCurrency(announcement.price)}
                  </td>
                  <td className="py-4 px-6 text-center">
                    <Link href={`/clonar/${announcement.id}`} className="button">
                      <button
                        type="button"
                        className="group relative flex w-10/12 justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                      >
                        <span className="absolute inset-y-0 left-0 flex items-center pl-3"></span>
                        Clonar
                      </button>
                    </Link>
                  </td>
                </tr>
              );
            })}
        </tbody>
      </table>
      {announcements && announcements.length === 0 && (
        <div className="overflow-hidden rounded-lg bg-slate-200 mt-2">
          <div className="px-4 py-5 sm:p-6 text-center">
            <svg
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              className="mx-auto h-20 w-20 text-gray-400"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                vectorEffect="non-scaling-stroke"
                strokeWidth={2}
                d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
              />
            </svg>
            <h3 className="mt-3 text-base font-semibold text-gray-900">
              Nada por aqui...
            </h3>
            <p className="mt-3 text-xs text-gray-500">
              Você não possui anúncios com esse SKU.
            </p>
            <div className="mt-3">
            <a
                  href={`https://www.mercadolivre.com.br/anuncios/lista?page=1&search=${sku}&sort=DEFAULT`}
                  target="_blank"
                  rel="noreferrer"
                >
              <button
                type="button"
                className="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                <ArrowTopRightOnSquareIcon
                  className="-ml-0.5 mr-1.5 h-5 w-5"
                  aria-hidden="true"
                />

                  Procurar no Mercado Livre
              </button>
              </a>{" "}

            </div>
          </div>
        </div>
      )}
    </div>
  );
}
