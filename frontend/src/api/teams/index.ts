import { api } from '../request'
import type { TeamListQuery, TeamRecord } from './types'

export function fetchTeams(query?: TeamListQuery) {
  return api.get<TeamRecord[]>('/teams/', query)
}

export type { TeamListQuery, TeamRecord } from './types'
