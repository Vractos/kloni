import { noLayoutRoutes, routePaths } from './constants/routes'
import LoginButton from './pages/login'
import LogoutButton from './pages/logout'
import { BrowserRouterProps } from 'react-router-dom'
import Pages from './pages'

const App: React.FC<BrowserRouterProps> = () => {


  return (
    <>
      <Pages/>
    </>
  )
}

export default App
