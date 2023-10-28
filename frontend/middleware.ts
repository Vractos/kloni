// middleware.js
import { withMiddlewareAuthRequired, getAccessToken } from '@auth0/nextjs-auth0/edge';
import { NextResponse } from 'next/server';

export default withMiddlewareAuthRequired(
  async function middleware(req) {
    const res = NextResponse.next();
    if (req.nextUrl.pathname.startsWith('/clonar')) {
      const { accessToken } = await getAccessToken(req, res);
      res.cookies.set('t', accessToken!);
    }
    return res;
  });