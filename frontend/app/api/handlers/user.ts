import 'server-only'

async function signingUpViaPurchaseAtPerfectPay(email: string, purchase_code: string): Promise<void> {
  const body = {
    email: email,
    purchase_code: purchase_code,
  }

  try {
    const res = await fetch(`${process.env.API_URL}/auth/perfectpay`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    })

    if (!res.ok) {
      throw new Error('Failed to clone announcement')
    }

  } catch (error) {
    console.log(error)
    throw new Error('Failed to clone announcement')
  }
}


export const runtime = 'edge';
export { signingUpViaPurchaseAtPerfectPay }