import { useState, useCallback } from 'react'
import { configService } from '../services/configService'
import type { Route } from '../apis/proto/opengate/v1/config'

export const useRoutes = () => {
  const [routes, setRoutes] = useState<Route[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadRoutes = useCallback(async () => {
    setLoading(true)
    setError(null)

    try {
      const response = await configService.getRoutes()
      setRoutes(response.routes || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load routes')
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    routes,
    loading,
    error,
    loadRoutes,
  }
}
