import { Container, Box, Typography, Card, CardContent, Skeleton, Alert } from '@mui/material'
import { useNavigate } from 'react-router-dom'
import RouteIcon from '@mui/icons-material/AltRoute'
import SettingsIcon from '@mui/icons-material/Settings'
import { useRoutes } from '../../hooks/useRoutes'
import { useEffect } from 'react'
import { PageHeader } from '../../components'

interface StatCardProps {
  title: string
  count: string | number
  icon: React.ReactNode
  color: string
  onClick?: () => void
}

const StatCard = ({ title, count, icon, color, onClick }: StatCardProps) => {
  return (
    <Card 
      sx={{ 
        borderRadius: 4,
        cursor: onClick ? 'pointer' : 'default',
        transition: 'all 0.3s',
        '&:hover': onClick ? {
          transform: 'translateY(-4px)',
          boxShadow: 6,
        } : {},
      }}
      onClick={onClick}
    >
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Box
            sx={{
              backgroundColor: `${color}15`,
              color: color,
              borderRadius: 3,
              width: 56,
              height: 56,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            {icon}
          </Box>
          <Box sx={{ flex: 1 }}>
            <Typography variant="body2" color="textSecondary" gutterBottom>
              {title}
            </Typography>
            <Typography variant="h4" sx={{ fontWeight: 600 }}>
              {typeof count === 'number' ? count.toLocaleString() : count}
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  )
}

const StatCardSkeleton = () => {
  return (
    <Card sx={{ borderRadius: 4 }}>
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Skeleton variant="rounded" width={56} height={56} />
          <Box sx={{ flex: 1 }}>
            <Skeleton variant="text" width="60%" />
            <Skeleton variant="text" width="40%" height={40} />
          </Box>
        </Box>
      </CardContent>
    </Card>
  )
}

export const DashboardPage = () => {
  const navigate = useNavigate()
  const { routes, loading, error, loadRoutes } = useRoutes()

  useEffect(() => {
    loadRoutes()
  }, [loadRoutes])

  const stats = [
    {
      title: 'Active Routes',
      count: routes.length,
      icon: <RouteIcon sx={{ fontSize: 28 }} />,
      color: '#2196f3',
      onClick: () => navigate('/gateway/configs'),
    },
    {
      title: 'Settings',
      count: '-',
      icon: <SettingsIcon sx={{ fontSize: 28 }} />,
      color: '#9c27b0',
      onClick: () => navigate('/gateway/settings'),
    },
  ]

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <PageHeader 
        title="Dashboard" 
        subtitle="Overview of your OpenGate API Gateway"
      />
      
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <Box
        sx={{
          display: 'grid',
          gridTemplateColumns: {
            xs: '1fr',
            sm: 'repeat(2, 1fr)',
            md: 'repeat(3, 1fr)',
          },
          gap: 3,
        }}
      >
        {loading
          ? Array.from({ length: 2 }).map((_, idx) => (
              <StatCardSkeleton key={idx} />
            ))
          : stats.map((stat) => (
              <StatCard
                key={stat.title}
                title={stat.title}
                count={stat.count}
                icon={stat.icon}
                color={stat.color}
                onClick={stat.onClick}
              />
            ))}
      </Box>

      {/* Routes Overview */}
      {!loading && routes.length > 0 && (
        <Box sx={{ mt: 4 }}>
          <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>
            Active Routes
          </Typography>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {routes.slice(0, 5).map((route) => (
              <Card key={route.name} sx={{ borderRadius: 2 }}>
                <CardContent sx={{ py: 2 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Box>
                      <Typography variant="subtitle1" fontWeight={600}>
                        {route.name}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {route.pathPrefix} → {route.targetUrl}
                      </Typography>
                    </Box>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      {route.stripPrefix && (
                        <Typography variant="caption" sx={{ 
                          bgcolor: 'info.light', 
                          color: 'info.contrastText',
                          px: 1, 
                          py: 0.5, 
                          borderRadius: 1 
                        }}>
                          Strip Prefix
                        </Typography>
                      )}
                      {route.authentication?.required && (
                        <Typography variant="caption" sx={{ 
                          bgcolor: 'warning.light', 
                          color: 'warning.contrastText',
                          px: 1, 
                          py: 0.5, 
                          borderRadius: 1 
                        }}>
                          Auth Required
                        </Typography>
                      )}
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            ))}
          </Box>
        </Box>
      )}
    </Container>
  )
}
