import { api } from '@/services/request'
import { apiBaseURL } from '..'
import type {
  CreateDeploymentPayload,
  Deployment,
  DeploymentListQuery,
  UpdateDeploymentPayload,
} from './types'

const baseURL = apiBaseURL.deployments

export function fetchProjectDeployments(projectId: string, query?: DeploymentListQuery) {
  return api.get<Deployment[]>(`${apiBaseURL.projects}/${projectId}/deployments`, query)
}

export function fetchDeploymentById(deploymentId: string) {
  return api.get<Deployment>(`${baseURL}/${deploymentId}`)
}

export function createDeployment(projectId: string, payload: CreateDeploymentPayload) {
  return api.post<Deployment>(`${apiBaseURL.projects}/${projectId}/deployments`, payload)
}

export function updateDeployment(deploymentId: string, payload: UpdateDeploymentPayload) {
  return api.patch<Deployment>(`${baseURL}/${deploymentId}`, payload)
}

export function deleteDeployment(deploymentId: string) {
  return api.delete<null>(`${baseURL}/${deploymentId}`)
}

export type {
  CreateDeploymentPayload,
  Deployment,
  DeploymentListQuery,
  UpdateDeploymentPayload,
} from './types'
