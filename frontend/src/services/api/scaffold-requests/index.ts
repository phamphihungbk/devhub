import { api } from '@/services/request'
import { apiBaseURL } from '../constants'
import type {
  CreateScaffoldRequestPayload,
  ScaffoldRequestSuggestion,
  ScaffoldRequestListQuery,
  ScaffoldRequestRecord,
  SuggestScaffoldRequestPayload,
} from './types'

const baseURL = apiBaseURL.scaffoldRequests

export function fetchProjectScaffoldRequests(projectId: string, query?: ScaffoldRequestListQuery) {
  return api.get<ScaffoldRequestRecord[]>(`${apiBaseURL.projects}/${projectId}/scaffold-requests`, query)
}

export function fetchScaffoldRequestById(scaffoldRequestId: string) {
  return api.get<ScaffoldRequestRecord>(`${baseURL}/${scaffoldRequestId}`)
}

export function createScaffoldRequest(projectId: string, payload: CreateScaffoldRequestPayload) {
  return api.post<ScaffoldRequestRecord>(`${apiBaseURL.projects}/${projectId}/scaffold-requests`, payload)
}

export function suggestProjectScaffoldRequest(projectId: string, payload: SuggestScaffoldRequestPayload) {
  return api.post<ScaffoldRequestSuggestion>(`${apiBaseURL.projects}/${projectId}/scaffold-suggestions`, payload)
}

export function deleteScaffoldRequest(scaffoldRequestId: string) {
  return api.delete<null>(`${baseURL}/${scaffoldRequestId}`)
}

export type {
  CreateScaffoldRequestPayload,
  ScaffoldRequestSuggestion,
  ScaffoldRequestListQuery,
  ScaffoldRequestRecord,
  ScaffoldRequestVariables,
  SuggestScaffoldRequestPayload,
} from './types'
