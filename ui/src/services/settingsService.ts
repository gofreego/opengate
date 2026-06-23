import { httpClient } from '../utils/httpClient'

const BASE_URL = '/opengate/v1'

export interface CORSConfig {
  enabled: boolean
  allowedOrigins: string
  allowedMethods: string
  allowedHeaders: string
  maxAge: number
}

export interface AppSettingsResponse {
  settings: Record<string, unknown>
  message: string
}

export const CORS_CONFIG_KEY = 'cors_config'

export const settingsService = {
  async getAll(): Promise<AppSettingsResponse> {
    const response = await httpClient.get<AppSettingsResponse>(`${BASE_URL}/app-settings`)
    return response.data
  },

  async upsert(key: string, value: unknown): Promise<void> {
    await httpClient.put(`${BASE_URL}/app-settings`, { key, value })
  },

  async getCORSConfig(): Promise<CORSConfig | null> {
    const resp = await settingsService.getAll()
    const raw = resp.settings?.[CORS_CONFIG_KEY]
    if (!raw) return null
    return raw as CORSConfig
  },

  async saveCORSConfig(config: CORSConfig): Promise<void> {
    await settingsService.upsert(CORS_CONFIG_KEY, config)
  },
}
