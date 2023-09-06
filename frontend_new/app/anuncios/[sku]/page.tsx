import { getAnnouncements } from '../../api/handlers/announcements'
import Table from '../../components/table'

export default async function TableAnnouncements({params}: {params: {sku: string}}) {
  const announcements = await getAnnouncements(params.sku, "test")
  
  return (
    <Table announcements={announcements} />
  )
}