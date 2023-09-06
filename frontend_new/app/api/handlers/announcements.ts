import { IAnnouncement } from '../../_lib/interfaces/announcements'

async function getAnnouncements(sku: string, token: string): Promise<IAnnouncement[]> {
  const res = await fetch(`${process.env.API_URL}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
    }
  })
 
  if (!res.ok) {
    // This will activate the closest `error.js` Error Boundary
    throw new Error('Failed to fetch data')
  }
 
  return res.json()
}

async function cloneAnnouncement(rootID: string,titles: string[],token: string): Promise<void> {
  const res = await fetch(`${process.env.API_URL}`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    }
  })
 
  if (!res.ok) {
    // This will activate the closest `error.js` Error Boundary
    throw new Error('Failed to fetch data')
  }
 
  return res.json()
}

export { getAnnouncements, cloneAnnouncement }
