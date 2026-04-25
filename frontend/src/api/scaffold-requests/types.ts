export interface ScaffoldRequestListQuery {
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface ScaffoldRequestVariables {
  service_name: string
  module_path: string
  port: number
  database: string
  enable_logging: boolean
}

export interface ScaffoldRequestRecord {
  id: string
  plugin_id: string
  requested_by: string
  project_id: string
  status: string
  environment?: string
  environments?: string
  variables: ScaffoldRequestVariables
}

export interface CreateScaffoldRequestPayload {
  plugin_id: string
  environment: string
  variables: ScaffoldRequestVariables
}

export interface SuggestScaffoldRequestPayload {
  prompt: string
  project_name?: string
  project_description?: string
  environment?: string
  environments?: string[]
}

export interface ScaffoldRequestSuggestion {
  source: string
  plugin_id: string
  plugin_name: string
  environment: string
  environments: string[]
  variables: ScaffoldRequestVariables
  rationale: string[]
}
