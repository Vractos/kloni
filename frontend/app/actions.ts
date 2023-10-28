"use server"

import { cloneAnnouncement } from '@/api/handlers/announcements'
import { getAccessToken } from '@auth0/nextjs-auth0'
import { revalidatePath } from 'next/cache'
import { permanentRedirect } from 'next/navigation'

export async function clone(preState: any, formData: FormData) {

  let success = false
  const titles = formData.getAll('title') as string[]
  const rootID = formData.get('id') as string
  const sku = formData.get('sku') as string

  try {
    await cloneAnnouncement(rootID, titles)
    success = true
  } catch (e) {
    return { fails: preState.fails += 1 };
  }
  finally {
    if (success) {
      revalidatePath('/?q=' + sku)
      permanentRedirect("../?q=" + sku)
    }
  }
}