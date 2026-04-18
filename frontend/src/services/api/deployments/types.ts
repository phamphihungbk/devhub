export interface DeploymentListQuery {
  startDate?: string
  endDate?: string
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface Deployment {
  id: string
  project_id: string
  plugin_id: string
  environment: string
  service: string
  version: string
  status: string
  external_ref?: string
  commit_sha?: string
  finished_at?: string
  triggered_by: string
}

export interface CreateDeploymentPayload {
  plugin_id: string
  environment: string
  service: string
  version: string
}

export interface UpdateDeploymentPayload {
  environment?: string
  service?: string
  version?: string
  status?: string
  external_ref?: string
  commit_sha?: string
  finished_at?: string
}
