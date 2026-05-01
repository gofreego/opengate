import { Container, Box, Typography, Card, CardContent, Skeleton, Alert } from '@mui/material'
import RouteIcon from '@mui/icons-material/AltRoute'
import { useStats } from '../../hooks/useStats'
import { useEffect } from 'react'
import { PageHeader } from '../../components'

interface StatCardProps {
  title: string
  count: string | number
  icon: React.ReactNode
  color: string
}

const StatCard = ({ title, count, icon, color }: StatCardProps) => {
  return (
    <Card 
      sx={{ 
        borderRadius: 4,
      }}
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
  const { stats, loading, error, loadStats } = useStats()

  useEffect(() => {
    loadStats()
  }, [loadStats])

  const statCards = [
    {
      title: 'Total Routes',
      count: stats?.totalRoutes ?? 0,
      icon: <RouteIcon sx={{ fontSize: 28 }} />,
      color: '#2196f3',
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
          ? Array.from({ length: 1 }).map((_, idx) => (
              <StatCardSkeleton key={idx} />
            ))
          : statCards.map((stat) => (
              <StatCard
                key={stat.title}
                title={stat.title}
                count={stat.count}
                icon={stat.icon}
                color={stat.color}
              />
            ))}
      </Box>
    </Container>
  )
}
