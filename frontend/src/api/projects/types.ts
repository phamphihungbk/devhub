export interface ProjectEnvironment {
  id: string
  name: string
  tier?: string
}

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
  environments?: Array<string | ProjectEnvironment>
  status?: string
  team_id?: string
  owner_team_id?: string
  owner_team_name?: string
  owner_team?: string
  scm_provider?: string
  owner_contact?: string
  created_by?: string
  created_by_name?: string
}

export interface ProjectPayload {
  name: string
  description?: string
  owner_team_id: string
}

export interface UpdateProjectPayload {
  name?: string
  description?: string
  environments?: Array<string | ProjectEnvironment>
  status?: string
  owner_team?: string
  scm_provider?: string
  owner_contact?: string
}
