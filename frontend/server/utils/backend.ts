import { createError } from 'h3'

interface BackendResponse<T> {
  code: string
  data?: T
  metadata?: unknown
  message?: string
}

export async function fetchBackend<T>(path: string): Promise<T> {
  const runtimeConfig = useRuntimeConfig()

  try {
    const response = await $fetch<BackendResponse<T>>(path, {
      baseURL: runtimeConfig.backendBaseUrl,
    })

    return response.data as T
  } catch (error) {
    throw createError({
      statusCode: 502,
      statusMessage: `Failed to fetch backend resource: ${path}`,
      data: error,
    })
  }
}
