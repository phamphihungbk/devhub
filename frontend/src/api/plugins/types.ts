export interface PluginListQuery {
  startDate?: string
  endDate?: string
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface PluginRecord {
  id: string
  name: string
  type: string
  version: string
  runtime: string
  entrypoint: string
  enabled?: boolean
  scope: string
  description: string
}

export interface PluginPayload {
  name: string
  version: string
  type: string
  runtime: string
  entrypoint: string
  scope: string
  description: string
  enabled?: boolean
}

export interface UpdatePluginPayload {
  name?: string
  description?: string
  type?: string
  version?: string
  runtime?: string
  entrypoint?: string
  scope?: string
  enabled?: boolean
}

export interface PluginSyncResult {
  discovered: number
  created: number
  updated: number
}
