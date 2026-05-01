import { ThemeProvider, SidebarLayout, NotificationProvider } from '@gofreego/tsutils'
import DashboardIcon from '@mui/icons-material/Dashboard'
import RouteIcon from '@mui/icons-material/AltRoute'
import { DashboardPage } from './pages/dashboard/DashboardPage'
import { ConfigsPage } from './pages/configs/ConfigsPage'

function App() {
  const menuItems = [
    {
      id: 'dashboard',
      label: 'Dashboard',
      path: '/gateway/dashboard',
      icon: <DashboardIcon />,
      component: <DashboardPage />,
    },
    {
      id: 'routes',
      label: 'Routes',
      path: '/gateway/routes',
      icon: <RouteIcon />,
      component: <ConfigsPage />,
    },
  ]

  return (
    <ThemeProvider>
      <NotificationProvider>
        <SidebarLayout
          menuItems={menuItems}
          defaultSelected="dashboard"
          isRouter={true}
          style={{ height: '100vh' }}
        />
      </NotificationProvider>
    </ThemeProvider>
  )
}

export default App
