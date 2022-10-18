import React from 'react'

interface IProductsData {
  picture: string,
  title: string,
  price: number,
  sku: string
}
interface ITableProps {
  products: IProductsData[]
}

const Table: React.FC<ITableProps> = ({ products }) => {
  return (

    <div className="overflow-x-auto relative rounded-lg">
      <table className="w-full text-sm text-left text-gray-500 table-auto">
        <thead className="text-xs text-gray-700 uppercase bg-gray-50">
          <tr>
            <th scope="col" className="py-3 px-6">
              Imagem
            </th>
            <th scope="col" className="py-3 px-6">
              Título do anúncio
            </th>
            <th scope="col" className="py-3 px-6">
              SKU
            </th>
            <th scope="col" className="py-3 px-6">
              Preço
            </th>
          </tr>
        </thead>
        <tbody>
          { products && products.map((product, index) => {
            return <tr className="bg-white border-b" key={index}>
              <th scope="row" className="py-4 px-6 font-medium text-gray-900 whitespace-nowrap">
                {product.picture}
              </th>
              <td className="py-4 px-6">
                {product.title}
              </td>
              <td className="py-4 px-6">
                {product.sku}
              </td>
              <td className="py-4 px-6">
                {product.price}
              </td>
            </tr>
          })}
        </tbody>
      </table>
    </div>

  )
}

export default Table