import { apiBaseURL } from '../constants'
import { api } from '../request'
import type { CreateUserPayload, UpdateUserPayload, UserListQuery, UserRecord } from './types'

const baseURL = `${apiBaseURL.users}/`

export function fetchUsers(query?: UserListQuery) {
  return api.get<UserRecord[]>(baseURL, query)
}

export function fetchUserById(userId: string) {
  return api.get<UserRecord>(`${apiBaseURL.users}/${userId}`)
}

export function createUser(payload: CreateUserPayload) {
  return api.post<UserRecord>(baseURL, payload)
}

export function updateUser(userId: string, payload: UpdateUserPayload) {
  return api.patch<UserRecord>(`${apiBaseURL.users}/${userId}`, payload)
}

export function deleteUser(userId: string) {
  return api.delete<null>(`${apiBaseURL.users}/${userId}`)
}

export type { CreateUserPayload, UpdateUserPayload, UserListQuery, UserRecord } from './types'
