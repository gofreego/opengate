import { useState, useCallback } from 'react'
import { useNotification } from '@gofreego/tsutils'
import { configService } from '../services/configService'
import type { GetStatsResponse } from '../apis/proto/opengate/v1/config'

export const useStats = () => {
  const [stats, setStats] = useState<GetStatsResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { showNotification } = useNotification()

  const loadStats = useCallback(async () => {
    setLoading(true)
    setError(null)

    try {
      const response = await configService.getStats()
      setStats(response)
    } catch (err) {
      const message = 'Failed to load stats'
      setError(message)
      showNotification(message, 'error')
    } finally {
      setLoading(false)
    }
  }, [showNotification])

  return {
    stats,
    loading,
    error,
    loadStats,
  }
}
