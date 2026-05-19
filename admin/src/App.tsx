import { useEffect } from 'react'
import { ThemeProvider, SidebarLayout, NotificationProvider, LoginCallbackPage, NotFoundPage, ProtectedRoute } from '@gofreego/tsutils'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import DashboardIcon from '@mui/icons-material/Dashboard'
import RouteIcon from '@mui/icons-material/AltRoute'
import { DashboardPage } from './pages/dashboard/DashboardPage'
import { RoutesPage } from './pages/routes/RoutesPage'
import { authService, sessionManager } from './services'

const LOGIN_URL = import.meta.env.VITE_LOGIN_URL as string

function App() {
  useEffect(() => {
    authService.initializeAuth()
  }, [])

  const handleLoginFailed = () => {
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
  ]

  return (
    <ThemeProvider>
      <NotificationProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/gateway/login-callback" element={<LoginCallbackPage authService={authService} onLoginFailed={handleLoginFailed} />} />
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
              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </NotificationProvider>
    </ThemeProvider>
  )
}

export default App
