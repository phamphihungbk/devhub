export interface UserListQuery {
  startDate?: string
  endDate?: string
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface UserRecord {
  id: string
  name: string
  email: string
  role: string
  team_id: string
}

export interface CreateUserPayload {
  name?: string
  email: string
  password: string
  role: string
}

export interface UpdateUserPayload {
  name?: string
  role?: string
}
