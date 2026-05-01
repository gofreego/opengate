import { HttpClient } from '@gofreego/tsutils'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Get dev headers from environment (for development/testing)
const getDevHeaders = (): Record<string, string> => {
  const headers: Record<string, string> = {}
  
  const devUserId = import.meta.env.VITE_DEV_USER_ID
  const devProfileId = import.meta.env.VITE_DEV_PROFILE_ID

  if (devUserId) {
    headers['X-User-Id'] = devUserId
  }
  if (devProfileId) {
    headers['X-Profile-Id'] = devProfileId
  }

  return headers
}

// Create and export a configured HTTP client instance
export const httpClient = new HttpClient({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: getDevHeaders(),
})

export default httpClient
