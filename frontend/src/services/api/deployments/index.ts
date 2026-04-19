import { api } from '@/services/request'
import { apiBaseURL } from '../constants'
import type {
  CreateDeploymentPayload,
  Deployment,
  DeploymentListQuery,
  UpdateDeploymentPayload,
} from './types'

const baseURL = apiBaseURL.deployments

export function fetchServiceDeployments(serviceId: string, query?: DeploymentListQuery) {
  return api.get<Deployment[]>(`${apiBaseURL.services}/${serviceId}/deployments`, query)
}

export function fetchDeploymentById(deploymentId: string) {
  return api.get<Deployment>(`${baseURL}/${deploymentId}`)
}

export function createDeployment(serviceId: string, payload: CreateDeploymentPayload) {
  return api.post<Deployment>(`${apiBaseURL.services}/${serviceId}/deployments`, payload)
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
