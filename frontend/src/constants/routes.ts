export interface IRoutes {
  HOME: string
  LOGIN: string
}

export type RouteKeyTypes = keyof IRoutes

export const routePaths: IRoutes = {
  HOME: '/',
  LOGIN: '/login'
}

export const noLayoutRoutes = ['LOGIN'] as const;