export interface ProjectListQuery {
  startDate?: string
  endDate?: string
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface Project {
  id: string
  name: string
  description: string
  environments: string[]
  status?: string
  owner_team?: string
  repo_url?: string
  repo_provider?: string
  owner_contact?: string
  created_by: string
}

export interface ProjectPayload {
  name: string
  description?: string
  environments: string[]
  status: string
  owner_team: string
  repo_url: string
  repo_provider: string
  owner_contact: string
}

export interface UpdateProjectPayload {
  name?: string
  description?: string
  environments?: string[]
  status?: string
  owner_team?: string
  repo_url?: string
  repo_provider?: string
  owner_contact?: string
}
