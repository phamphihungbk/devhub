import { api } from '@/services/request'
import { apiBaseURL } from '..'
import type { Project, ProjectListQuery, ProjectPayload, UpdateProjectPayload } from './types'

const baseURL = `${apiBaseURL.projects}/`

export function fetchProjects(query?: ProjectListQuery) {
  return api.get<Project[]>(baseURL, query)
}

export function fetchProjectById(projectId: string) {
  return api.get<Project>(`${apiBaseURL.projects}/${projectId}`)
}

export function createProject(payload: ProjectPayload) {
  return api.post<Project>(baseURL, payload)
}

export function updateProject(projectId: string, payload: UpdateProjectPayload) {
  return api.patch<Project>(`${apiBaseURL.projects}/${projectId}`, payload)
}

export function deleteProject(projectId: string) {
  return api.delete<null>(`${apiBaseURL.projects}/${projectId}`)
}

export type { Project, ProjectListQuery, ProjectPayload, UpdateProjectPayload } from './types'
