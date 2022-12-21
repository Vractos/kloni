export interface IRoutes {
  HOME: string
  LOGIN: string
  NOT_FOUND: string
}

export type RouteKeyTypes = keyof IRoutes

export const routePaths: IRoutes = {
  HOME: '/',
  LOGIN: '/login',
  NOT_FOUND: '/404'
}
export const routeNames: IRoutes = {
  HOME: 'Inicio',
  LOGIN: 'Login',
  NOT_FOUND: 'Not Found'
}

export const navBarRoutes = ['HOME'] as const;

export const noLayoutRoutes = ['LOGIN'] as const;