import { api } from '@/services/request'
import { apiBaseURL } from '..'
import type { Credentials, LoginTokens, UserProfile } from './types'

export async function loginWithPassword(credentials: Credentials): Promise<LoginTokens> {
  const response = await api.post<{ access_token: string; refresh_token: string }>(`${apiBaseURL.auth}/login`, credentials)

  return {
    accessToken: response.access_token,
    refreshToken: response.refresh_token,
  }
}

export function logoutCurrentUser() {
  return api.post<null>(`${apiBaseURL.auth}/logout`)
}

export function fetchCurrentUser() {
  return api.get<UserProfile>(`${apiBaseURL.auth}/me`)
}

export type { Credentials, LoginTokens, UserProfile } from './types'
