import { api } from '@/services/request'
import { apiBaseURL } from '../constants'
import type { ScaffoldRequestVariables } from '../scaffold-requests/types'

export interface Service {
  id: string
  project_id: string
  name: string
  repo_url: string
}

export interface SuggestScaffoldPayload {
  service_name?: string
  project_name?: string
  project_description?: string
  repo_url?: string
  environment?: string
  environments?: string[]
}

export interface ScaffoldSuggestion {
  source: string
  environment: string
  variables: ScaffoldRequestVariables
  rationale: string[]
}

export function fetchProjectServices(projectId: string) {
  return api.get<Service[]>(`${apiBaseURL.projects}/${projectId}/services`)
}

export function suggestServiceScaffold(serviceId: string, payload: SuggestScaffoldPayload) {
  return api.post<ScaffoldSuggestion>(`${apiBaseURL.services}/${serviceId}/scaffold-suggestions`, payload)
}
