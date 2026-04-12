import { requestApi } from '../index'
import { toQueryParams } from '../util'
import type {
  BackendUpdateUserResponse,
  BackendUser,
  CreateUserInput,
  FindAllUsersQuery,
  PaginationMetadata,
  UpdateUserInput,
  User,
} from './types'
import { mapBackendUpdatedUser, mapBackendUser, mapBackendUsers } from './types'

export async function findAllUsers(query?: FindAllUsersQuery) {
  const response = await requestApi<BackendUser[], PaginationMetadata>('/users/', {
    method: 'GET',
    query: toQueryParams(query),
  })

  return {
    ...response,
    data: response.data ? mapBackendUsers(response.data) : [],
  }
}

export async function findUserById(userId: string) {
  const response = await requestApi<BackendUser>(`/users/${userId}`, {
    method: 'GET',
  })

  return {
    ...response,
    data: response.data ? mapBackendUser(response.data) : undefined,
  }
}

export async function createUser(input: CreateUserInput) {
  const response = await requestApi<BackendUser>(`/users/`, {
    method: 'POST',
    body: input,
  })

  return {
    ...response,
    data: response.data ? mapBackendUser(response.data) : undefined,
  }
}

export async function updateUser(userId: string, input: UpdateUserInput) {
  const response = await requestApi<BackendUpdateUserResponse>(`/users/${userId}`, {
    method: 'PATCH',
    body: input,
  })

  return {
    ...response,
    data: response.data ? mapBackendUpdatedUser(response.data) : undefined,
  }
}

export async function deleteUser(userId: string) {
  return requestApi<null>(`/users/${userId}`, {
    method: 'DELETE',
  })
}
