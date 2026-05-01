import { useState, useCallback } from 'react'
import { useNotification } from '@gofreego/tsutils'
import { configService } from '../services/configService'
import type { Config, CreateConfigRequest, UpdateConfigRequest } from '../apis/proto/opengate/v1/config'

interface UseConfigsOptions {
  limit?: number
  offset?: number
  search?: string
  append?: boolean
}

export const useConfigs = () => {
  const [configs, setConfigs] = useState<Config[]>([])
  const [selectedConfig, setSelectedConfig] = useState<Config | null>(null)
  const [loading, setLoading] = useState(false)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [total, setTotal] = useState(0)
  const [hasMore, setHasMore] = useState(false)
  const { showNotification } = useNotification()

  const loadConfigs = useCallback(async (options: UseConfigsOptions = {}) => {
    const { limit = 10, offset = 0, search, append = false } = options
    
    if (append) {
      setLoadingMore(true)
    } else {
      setLoading(true)
    }
    setError(null)

    try {
      const response = await configService.list({ limit, offset, search: search || '' })
      const newConfigs = response.configs || []
      
      if (append) {
        setConfigs(prev => [...prev, ...newConfigs])
      } else {
        setConfigs(newConfigs)
      }
      
      setTotal(response.total || 0)
      setHasMore(offset + newConfigs.length < (response.total || 0))
    } catch (err) {
      const message = 'Failed to load routes'
      setError(message)
      showNotification(message, 'error')
    } finally {
      setLoading(false)
      setLoadingMore(false)
    }
  }, [showNotification])

  const getConfigById = useCallback(async (id: string): Promise<Config | null> => {
    try {
      const response = await configService.getById(id)
      setSelectedConfig(response.config || null)
      return response.config || null
    } catch (err) {
      const message = 'Failed to load route'
      setError(message)
      showNotification(message, 'error')
      return null
    }
  }, [showNotification])

  const createConfig = useCallback(async (data: CreateConfigRequest): Promise<Config | null> => {
    setLoading(true)
    setError(null)

    try {
      const response = await configService.create(data)
      showNotification('Route created successfully', 'success')
      return response.config || null
    } catch (err) {
      const message = 'Failed to create route'
      setError(message)
      showNotification(message, 'error')
      return null
    } finally {
      setLoading(false)
    }
  }, [showNotification])

  const updateConfig = useCallback(async (id: string, data: Omit<UpdateConfigRequest, 'id'>): Promise<Config | null> => {
    setLoading(true)
    setError(null)

    try {
      const response = await configService.update(id, data)
      showNotification('Route updated successfully', 'success')
      return response.config || null
    } catch (err) {
      const message = 'Failed to update route'
      setError(message)
      showNotification(message, 'error')
      return null
    } finally {
      setLoading(false)
    }
  }, [showNotification])

  const deleteConfig = useCallback(async (id: string): Promise<boolean> => {
    setLoading(true)
    setError(null)

    try {
      await configService.delete(id)
      setConfigs(prev => prev.filter(c => c.id !== id))
      showNotification('Route deleted successfully', 'success')
      return true
    } catch (err) {
      const message = 'Failed to delete route'
      setError(message)
      showNotification(message, 'error')
      return false
    } finally {
      setLoading(false)
    }
  }, [showNotification])

  const clearSelectedConfig = useCallback(() => {
    setSelectedConfig(null)
  }, [])

  return {
    configs,
    selectedConfig,
    loading,
    loadingMore,
    error,
    total,
    hasMore,
    loadConfigs,
    getConfigById,
    createConfig,
    updateConfig,
    deleteConfig,
    clearSelectedConfig,
  }
}
