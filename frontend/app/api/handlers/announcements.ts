import 'server-only';

import { IAnnouncement } from '../../../lib/interfaces/announcements'
import { getAccessToken } from '@auth0/nextjs-auth0/edge';

async function getAnnouncements(sku: string): Promise<IAnnouncement[]> {
  const { accessToken } = await getAccessToken()

  try {
    const res = await fetch(`${process.env.API_URL}/announcement/${sku}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      }
    })

    if (!res.ok) {
      if (res.status === 404) {
        return [];
      } else {
        throw new Error('Failed to fetch data. Server error')
      }
    }

    const data = await res.json() as IAnnouncement[]

    if (!data || data.length === 0) {
      console.log('Empty JSON response');
      return [];
    }

    data.map((announcement: IAnnouncement): IAnnouncement => ({
      id: announcement.id,
      title: announcement.title,
      picture: announcement.picture,
      price: announcement.price,
      quantity: announcement.quantity,
      sku: announcement.sku,
      link: announcement.link
    }))

    return data
  } catch (error) {
    console.log(error)
    throw new Error('Failed to fetch data')
  }
}

async function cloneAnnouncement(rootID: string, titles: string[], accessToken: string): Promise<void> {
  const body = {
    root_id: rootID,
    titles: titles,
  }

  try {
    const res = await fetch(`${process.env.API_URL}/announcement`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
      body: JSON.stringify(body)
    })

    if (!res.ok) {
      throw new Error('Failed to clone announcement')
    }
  } catch (error) {
    throw new Error('Failed to clone announcement')
  }
}

export const runtime = 'edge';
export { getAnnouncements, cloneAnnouncement }
