import { httpClient } from '../utils/httpClient'
import type {
  Config,
  CreateConfigRequest,
  CreateConfigResponse,
  GetConfigResponse,
  ListConfigsRequest,
  ListConfigsResponse,
  UpdateConfigRequest,
  UpdateConfigResponse,
  DeleteConfigResponse,
  GetRoutesResponse,
  GetStatsResponse,
} from '../apis/proto/opengate/v1/config'

const BASE_URL = '/opengate/v1'

export const configService = {
  async list(params: ListConfigsRequest): Promise<ListConfigsResponse> {
    const queryParams = new URLSearchParams({
      limit: params.limit.toString(),
      offset: params.offset.toString(),
      ...(params.search && { search: params.search }),
    })
    
    const response = await httpClient.get<ListConfigsResponse>(
      `${BASE_URL}/configs?${queryParams.toString()}`
    )
    return response.data
  },

  async getById(id: string): Promise<GetConfigResponse> {
    const response = await httpClient.get<GetConfigResponse>(`${BASE_URL}/configs/${id}`)
    return response.data
  },

  async create(data: CreateConfigRequest): Promise<CreateConfigResponse> {
    const response = await httpClient.post<CreateConfigResponse>(`${BASE_URL}/configs`, data)
    return response.data
  },

  async update(id: string, data: Omit<UpdateConfigRequest, 'id'>): Promise<UpdateConfigResponse> {
    const response = await httpClient.put<UpdateConfigResponse>(`${BASE_URL}/configs/${id}`, data)
    return response.data
  },

  async delete(id: string): Promise<DeleteConfigResponse> {
    const response = await httpClient.delete<DeleteConfigResponse>(`${BASE_URL}/configs/${id}`)
    return response.data
  },

  async getRoutes(): Promise<GetRoutesResponse> {
    const response = await httpClient.get<GetRoutesResponse>(`${BASE_URL}/routes`)
    return response.data
  },

  async getStats(): Promise<GetStatsResponse> {
    const response = await httpClient.get<GetStatsResponse>(`${BASE_URL}/stats`)
    return response.data
  },
}

// Helper to convert form data to API request
export const toCreateConfigRequest = (data: Partial<Config>): CreateConfigRequest => ({
  name: data.name || '',
  pathPrefix: data.pathPrefix || '',
  targetUrl: data.targetUrl || '',
  stripPrefix: data.stripPrefix || false,
  authentication: data.authentication,
  middleware: data.middleware || [],
  timeout: data.timeout || '30000000000',
})

export const toUpdateConfigRequest = (data: Config): UpdateConfigRequest => ({
  id: data.id,
  name: data.name,
  pathPrefix: data.pathPrefix,
  targetUrl: data.targetUrl,
  stripPrefix: data.stripPrefix,
  authentication: data.authentication,
  middleware: data.middleware || [],
  timeout: data.timeout || '30000000000',
})
