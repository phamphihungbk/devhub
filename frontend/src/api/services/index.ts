import { apiBaseURL } from '../constants'
import { api } from '../request'

export interface Service {
  id: string
  project_id: string
  name: string
  repo_url: string
}

export interface ServiceDependency {
  id: string
  service_id: string
  depends_on_service_id: string
  depends_on_service?: Service
  type: string
  protocol?: string
  port?: number | null
  path?: string
  config?: Record<string, unknown>
  created_by: string
}

export interface CreateServiceDependencyPayload {
  depends_on_service_id: string
  type: string
  protocol?: string
  port?: number | null
  path?: string
  config?: Record<string, unknown>
}

export function fetchProjectServices(projectId: string) {
  return api.get<Service[]>(`${apiBaseURL.projects}/${projectId}/services`)
}

export function fetchServiceDependencies(serviceId: string) {
  return api.get<ServiceDependency[]>(`${apiBaseURL.services}/${serviceId}/dependencies`)
}

export function createServiceDependency(serviceId: string, payload: CreateServiceDependencyPayload) {
  return api.post<ServiceDependency>(`${apiBaseURL.services}/${serviceId}/dependencies`, payload)
}

export function deleteServiceDependency(serviceId: string, dependencyId: string) {
  return api.delete<ServiceDependency>(`${apiBaseURL.services}/${serviceId}/dependencies/${dependencyId}`)
}
