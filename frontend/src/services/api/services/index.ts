import { api } from '@/services/request'
import { apiBaseURL } from '../constants'

export interface Service {
  id: string
  project_id: string
  name: string
  repo_url: string
}

export function fetchProjectServices(projectId: string) {
  return api.get<Service[]>(`${apiBaseURL.projects}/${projectId}/services`)
}
