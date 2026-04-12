type ApiSuccessResponse<TData, TMetadata = unknown> = {
  code: string
  data?: TData
  metadata?: TMetadata
}

type ApiErrorResponse = {
  code?: string
  message?: string
  data?: unknown
}

type RequestApiOptions = Omit<Parameters<typeof $fetch>[1], 'baseURL'>

function resolveApiBaseUrl() {
  const runtimeConfig = useRuntimeConfig()

  return import.meta.server
    ? runtimeConfig.backendBaseUrl
    : runtimeConfig.public.apiBaseUrl
}

export async function requestApi<TData, TMetadata = unknown>(
  path: string,
  options?: RequestApiOptions,
): Promise<ApiSuccessResponse<TData, TMetadata>> {
  const baseURL = resolveApiBaseUrl()

  try {
    return await $fetch<ApiSuccessResponse<TData, TMetadata>>(path, {
      baseURL,
      ...options,
    })
  } catch (error) {
    const apiError = error as { data?: ApiErrorResponse; statusMessage?: string; message?: string }
    const message = apiError.data?.message || apiError.statusMessage || apiError.message || 'API request failed'

    throw new Error(message)
  }
}
