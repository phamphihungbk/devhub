export interface ScaffoldRequestListQuery {
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface ScaffoldRequestVariables {
  service_name: string
  port: number
  database: string
  enable_logging: boolean
}

export interface ScaffoldRequestRecord {
  id: string
  plugin_id: string
  requested_by: string
  project_id: string
  template: string
  status: string
  environment?: string
  environments?: string
  variables: ScaffoldRequestVariables
}

export interface CreateScaffoldRequestPayload {
  plugin_id: string
  template: string
  environment: string
  variables: ScaffoldRequestVariables
}
