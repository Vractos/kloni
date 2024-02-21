'use server'

import { cloneAnnouncement } from '@/api/handlers/announcements'
import { signingUpViaPurchaseAtPerfectPay } from '@/api/handlers/user'
import { revalidatePath } from 'next/cache'
import { cookies } from 'next/headers'
import { permanentRedirect } from 'next/navigation'

export async function clone(preState: any, formData: FormData) {
  const accessToken = cookies().get('t')?.value

  let success = false
  const titles = formData.getAll('title') as string[]
  const rootID = formData.get('id') as string
  const sku = formData.get('sku') as string

  try {
    await cloneAnnouncement(rootID, titles, accessToken!)
    success = true
    // Block execution for 1 second to prevent abuse
    await new Promise(resolve => setTimeout(resolve, 5000))
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

export async function submitPerfectPayPurchase(preState: any, formData: FormData) {
  let success = false
  const email = formData.get('email') as string
  const purchaseCode = formData.get('purchase_code') as string

  try {
    await signingUpViaPurchaseAtPerfectPay(email, purchaseCode)
    success = true
    // Block execution for 1 second to prevent abuse
    await new Promise(resolve => setTimeout(resolve, 5000))
  } catch (e) {
    return { fails: preState.fails += 1 };
  }
  // finally {
  //   if (success) {
        

  //   }
  // }
}