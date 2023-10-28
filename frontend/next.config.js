/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '',
  images: {
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'http2.mlstatic.com',
        port: '',
      },
      {
        protocol: 'https',
        hostname: 'http2.mlstatic.com',
        port: '',
      },
      {
        protocol: 'https',
        hostname: 'tailwindui.com',
        port: '',
      },
    ],
  },
  experimental: {
    esmExternals: 'loose'
  },
}

module.exports = nextConfig
