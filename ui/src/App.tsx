import { useEffect, useState } from 'react'
import { ThemeProvider, SidebarLayout, NotificationProvider, LoginCallbackPage, ProtectedRoute } from '@gofreego/tsutils'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import DashboardIcon from '@mui/icons-material/Dashboard'
import RouteIcon from '@mui/icons-material/AltRoute'
import SettingsIcon from '@mui/icons-material/Settings'
import { DashboardPage } from './pages/dashboard/DashboardPage'
import { RoutesPage } from './pages/routes/RoutesPage'
import { SettingsPage } from './pages/settings'
import { authService, sessionManager } from './services'

const LOGIN_URL = import.meta.env.VITE_LOGIN_URL as string

function App() {
  const [isInitialized, setIsInitialized] = useState(false);
  useEffect(() => {
    authService.initializeAuth()
    setIsInitialized(true);
  }, [])

  const handleLoginFailed = () => {
    console.log("Login failed, redirecting to -> ", LOGIN_URL)
    window.location.href = LOGIN_URL
  }

  const menuItems = [
    {
      id: 'dashboard',
      label: 'Dashboard',
      path: '/gateway/dashboard',
      icon: <DashboardIcon />,
    },
    {
      id: 'routes',
      label: 'Routes',
      path: '/gateway/routes',
      icon: <RouteIcon />,
    },
    {
      id: 'settings',
      label: 'Settings',
      path: '/gateway/settings',
      icon: <SettingsIcon />,
    },
  ]

  if (!isInitialized) {
    // return loading spinner or placeholder
    return <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <div>Loading...</div>
    </div>
  }

  return (
    <ThemeProvider>
      <NotificationProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/gateway/login-callback" element={<LoginCallbackPage authService={authService} navigateTo="/gateway/dashboard" onLoginFailed={handleLoginFailed} />} />
            <Route
              path="/"
              element={
                <ProtectedRoute sessionManager={sessionManager} loginUrl={LOGIN_URL} callbackPath="/gateway/login-callback">
                  <SidebarLayout menuItems={menuItems} isRouter={true} isBrowserRouter={false} style={{ height: '100vh' }} />
                </ProtectedRoute>
              }
            >
              <Route path="gateway/dashboard" element={<DashboardPage />} />
              <Route path="gateway/routes" element={<RoutesPage />} />
              <Route path="gateway/settings" element={<SettingsPage />} />
              <Route path="*" element={<Navigate to="/gateway/dashboard" replace />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </NotificationProvider>
    </ThemeProvider>
  )
}

export default App
