import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { IAnnouncement } from '../../interfaces/announcement'
import { formatCurrency } from '../../utils/formatter'
import Modal from '../modal'

interface ITableProps {
  announcements: IAnnouncement[]
}

const Table: React.FC<ITableProps> = ({ announcements }) => {
  const [modalOpen, setModalOpen] = useState(false)

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
            <th scope="col" className="py-3 px-6 text-center">
            </th>
          </tr>
        </thead>
        <tbody>
          {announcements && announcements.map((announcement, index) => {
            return <tr className="bg-white border-b" key={index}>
              <th scope="row" className="py-4 px-6 font-medium text-gray-900 whitespace-nowrap">
                <div className="flex items-center justify-center">
                  <div className="w-10 h-10 flex-shrink-0">
                    <img className="rounded-full" src={announcement.picture} width="40" height="40" />
                  </div>
                </div>
              </th>
              <td className="py-4 px-6">
                <a href={announcement.link} target="_blank" rel="noreferrer" className='text-blue-600'>
                  {announcement.title}
                </a>
              </td>
              <td className="py-4 px-6">
                {announcement.sku}
              </td>
              <td className="py-4 px-6 text-center">
                {announcement.quantity}
              </td>
              <td className="py-4 px-6 text-center">
                {formatCurrency(announcement.price)}
              </td>
              <td className="py-4 px-6 text-center">
                <button
                  type="button"
                  className="group relative flex w-10/12 justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                  onClick={() => setModalOpen(true)}
>
                  <span className="absolute inset-y-0 left-0 flex items-center pl-3">
                  </span>
                  Clonar
                </button>
                <Modal isOpen={modalOpen} announcementId={announcement.id} handleClose={setModalOpen}/>
              </td>
            </tr>
          })}
        </tbody>
      </table>
    </div>

  )
}

export default Table