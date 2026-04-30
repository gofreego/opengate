// Types matching the proto definitions
export interface Authentication {
  type: string
  config: Record<string, string>
}

export interface Middleware {
  type: string
  config: Record<string, string>
}

export interface Config {
  id: string
  name: string
  pathPrefix: string
  targetUrl: string
  stripPrefix: boolean
  authentication?: Authentication
  middleware: Middleware[]
  timeout: number
  createdAt: string
  updatedAt: string
}

export interface Route {
  pathPrefix: string
  targetUrl: string
  stripPrefix: boolean
  authentication?: Authentication
  timeout: number
}

// Request/Response types
export interface CreateConfigRequest {
  name: string
  pathPrefix: string
  targetUrl: string
  stripPrefix: boolean
  authentication?: Authentication
  middleware: Middleware[]
  timeout: number
}

export interface CreateConfigResponse {
  config: Config
}

export interface GetConfigRequest {
  id: string
}

export interface GetConfigResponse {
  config: Config
}

export interface ListConfigsRequest {
  limit: number
  offset: number
  search?: string
}

export interface ListConfigsResponse {
  configs: Config[]
  total: number
}

export interface UpdateConfigRequest {
  id: string
  name: string
  pathPrefix: string
  targetUrl: string
  stripPrefix: boolean
  authentication?: Authentication
  middleware: Middleware[]
  timeout: number
}

export interface UpdateConfigResponse {
  config: Config
}

export interface DeleteConfigRequest {
  id: string
}

export interface DeleteConfigResponse {
  success: boolean
}

export interface GetRoutesRequest {}

export interface GetRoutesResponse {
  routes: Route[]
}
