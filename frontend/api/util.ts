export function toQueryParams<T extends Record<string, unknown>>(query?: T) {
  if (!query) {
    return undefined
  }

  return Object.fromEntries(
    Object.entries(query).filter(([, value]) => value !== undefined && value !== null),
  )
}
