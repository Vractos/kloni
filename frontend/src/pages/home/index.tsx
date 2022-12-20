import React, { FormEvent, useState } from 'react'
import NavBar from '../../components/layouts/navBar'
import { MagnifyingGlassIcon } from "@heroicons/react/24/solid";
import Table from '../../components/table';
import { IAnnouncement } from '../../interfaces/announcement';
import api from '../../api';
import { useAuth0 } from '@auth0/auth0-react';

function Home() {
  const [announcements, setAnnouncements] = useState<IAnnouncement[]>([])
  const [searchBarValue, setSearchBarValue] = useState("")
  const [emptyResultText, setEmptyResultText] = useState("Busque por anúncios através do SKU")
  const { getAccessTokenSilently } = useAuth0()


  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    try {
      const token = await getAccessTokenSilently()
      const anns = await api.announcement.getAnnouncements(searchBarValue, token)
      setAnnouncements(anns)
    } catch (error) {
      setEmptyResultText("Não há anúncios com esse SKU")
      setAnnouncements([])}
  }

  return (
    <div className="min-h-full">
      <NavBar />
      <header className="bg-white shadow">
        <div className="mx-auto max-w-7xl py-6 px-4 sm:px-6 lg:px-8">
          {/* <h1 className="text-3xl font-bold tracking-tight text-gray-900">Dashboard</h1> */}
          <form onSubmit={handleSubmit}>
            <label htmlFor="search" className="mb-2 text-sm font-medium sr-only dark:text-gray-300">Buscar</label>
            <div className="relative w-8/12 ">
              <div className="flex absolute inset-y-0 left-0 items-center pl-3 pointer-events-none">
                <MagnifyingGlassIcon className='h-5 w-5 text-indigo-500 group-hover:text-indigo-400' />
              </div>
              <input value={searchBarValue} type="search" id="search" className="outline-none block p-4 pl-10 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500" placeholder="SKU" onChange={e => setSearchBarValue(e.target.value)} required />
              <button type="submit" onSubmit={e => handleSubmit(e)} className="text-white absolute right-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2">Buscar</button>
            </div>
          </form>
        </div>
      </header>
      <main>
        <div className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
          <div className="px-4 py-6 sm:px-0 ">
            {announcements.length ?
              <Table announcements={announcements} />
              :
              <div className="px-4 py-6 sm:px-0">
                <p className="mt-5 text-center">
                  {emptyResultText}
                </p>
              </div>
            }
          </div>
        </div>
      </main>
    </div>
  )
}

export default Home;