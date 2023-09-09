import 'server-only';

import { IAnnouncement } from '../../_lib/interfaces/announcements'
import announcements from '../../mocks/announcements.json'


async function getAnnouncements(sku: string, token: string): Promise<IAnnouncement[]> {
  const retrievedAnnouncements = announcements.filter(announcements => announcements.sku.toLowerCase() === sku)

  return new Promise((res, rej) => {
    setTimeout(() => {
      if (!retrievedAnnouncements.length) {
        res([])
        // rej({
        //   message: "No announcements found for this sku"
        // });
      } else {
        res(retrievedAnnouncements);
      }
    }, 1000)
  });

  // const res = await fetch(`${process.env.API_URL}`, {
  //   method: 'GET',
  //   headers: {
  //     'Authorization': `Bearer ${token}`,
  //   }
  // })

  // if (!res.ok) {
  //   // This will activate the closest `error.js` Error Boundary
  //   throw new Error('Failed to fetch data')
  // }

  // return res.json()
}

async function cloneAnnouncement(rootID: string, titles: string[], token: string): Promise<void> {
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
