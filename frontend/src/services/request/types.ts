export interface ApiEnvelope<T> {
  code: string
  data: T
  message?: string
  metadata?: ApiMetadata
}

export interface PaginationMetadata {
  total_records?: number
  limit?: number
  offset?: number
}

export interface ApiMetadata {
  pagination?: PaginationMetadata
}

export type QueryValue = string | number | boolean | null | undefined
