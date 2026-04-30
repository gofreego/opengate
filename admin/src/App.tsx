import { ThemeProvider, SidebarLayout, NotificationProvider } from '@gofreego/tsutils'
import DashboardIcon from '@mui/icons-material/Dashboard'
import SettingsIcon from '@mui/icons-material/Settings'
import RouteIcon from '@mui/icons-material/AltRoute'
import { DashboardPage } from './pages/dashboard/DashboardPage'
import { ConfigsPage } from './pages/configs/ConfigsPage'
import { SettingsPage } from './pages/settings/SettingsPage'

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
      id: 'configs',
      label: 'Configs',
      path: '/gateway/configs',
      icon: <RouteIcon />,
      component: <ConfigsPage />,
    },
    {
      id: 'settings',
      label: 'Settings',
      path: '/gateway/settings',
      icon: <SettingsIcon />,
      component: <SettingsPage />,
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
