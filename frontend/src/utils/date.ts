export const formatRequestedAt = (value: string) => {
  return new Date(value).toLocaleString()
}

export const formatOptionalDate = (value?: string | null, fallback = 'Not set') => {
  return value ? new Date(value).toLocaleString() : fallback
}
