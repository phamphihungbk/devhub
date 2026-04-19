import { TOKEN_STORAGE_KEY } from '@/stores/modules/auth'
import type { ApiEnvelope, QueryValue } from './types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

function withQuery(path: string, query?: Record<string, QueryValue>) {
  if (!query) return path

  const searchParams = new URLSearchParams()

  for (const [key, value] of Object.entries(query)) {
    if (value === undefined || value === null || value === '') continue
    searchParams.set(key, String(value))
  }

  const search = searchParams.toString()
  return search ? `${path}?${search}` : path
}

async function request<T>(path: string, init: RequestInit = {}) {
  const token = localStorage.getItem(TOKEN_STORAGE_KEY)
  const headers = new Headers(init.headers)

  headers.set('Accept', 'application/json')

  if (init.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers,
  })

  const contentType = response.headers.get('content-type') || ''
  const payload = contentType.includes('application/json')
    ? ((await response.json()) as Partial<ApiEnvelope<T>> & { message?: string })
    : null

  if (!response.ok) {
    const fallback = response.status === 401 ? 'Your session is no longer valid.' : 'Request failed.'
    const message = payload?.message || fallback
    if (response.status === 401) {
      localStorage.removeItem(TOKEN_STORAGE_KEY)
    }
    throw new ApiError(message, response.status)
  }

  return (payload?.data as T) ?? (payload as T)
}

export const api = {
  get: <T>(path: string, query?: Record<string, QueryValue>) => request<T>(withQuery(path, query)),
  post: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    }),
  patch: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    }),
  delete: <T>(path: string) =>
    request<T>(path, {
      method: 'DELETE',
    }),
}
