export interface TeamListQuery {
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface TeamRecord {
  id: string
  name: string
  owner_contact: string
}
