import { api } from '@/services/request'
import { apiBaseURL } from '../constants'
import type { CreateReleasePayload, Release } from './types'

export function createRelease(serviceId: string, payload: CreateReleasePayload) {
  return api.post<Release>(`/services/${serviceId}/releases`, payload)
}

export type { CreateReleasePayload, Release } from './types'
