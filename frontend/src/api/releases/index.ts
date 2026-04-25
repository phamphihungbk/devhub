import { api } from '../request'
import type { CreateReleasePayload, Release } from './types'

export function fetchServiceReleases(serviceId: string) {
  return api.get<Release[]>(`/services/${serviceId}/releases`)
}

export function createRelease(serviceId: string, payload: CreateReleasePayload) {
  return api.post<Release>(`/services/${serviceId}/releases`, payload)
}

export type { CreateReleasePayload, Release } from './types'
