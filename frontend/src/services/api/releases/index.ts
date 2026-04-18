import { api } from '@/services/request'
import { apiBaseURL } from '../constants'
import type { CreateReleasePayload, Release } from './types'

export function createRelease(projectId: string, payload: CreateReleasePayload) {
  return api.post<Release>(`${apiBaseURL.projects}/${projectId}/releases`, payload)
}

export type { CreateReleasePayload, Release } from './types'
