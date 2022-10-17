export interface IRoutes {
  HOME: string
  LOGIN: string
}

export type RouteKeyTypes = keyof IRoutes

export const routePaths: IRoutes = {
  HOME: '/',
  LOGIN: '/login'
}
export const routeNames: IRoutes = {
  HOME: 'Inicio',
  LOGIN: 'Login'
}

export const navBarRoutes = ['HOME'] as const;

export const noLayoutRoutes = ['LOGIN'] as const;