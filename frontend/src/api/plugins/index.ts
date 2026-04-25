import { apiBaseURL } from '../constants'
import { api } from '../request'
import type { PluginListQuery, PluginPayload, PluginRecord, UpdatePluginPayload } from './types'

const baseURL = `${apiBaseURL.plugins}/`

export function fetchPlugins(query?: PluginListQuery) {
  return api.get<PluginRecord[]>(baseURL, query)
}

export function fetchPluginById(pluginId: string) {
  return api.get<PluginRecord>(`${apiBaseURL.plugins}/${pluginId}`)
}

export function createPlugin(payload: PluginPayload) {
  return api.post<PluginRecord>(baseURL, payload)
}

export function updatePlugin(pluginId: string, payload: UpdatePluginPayload) {
  return api.patch<PluginRecord>(`${apiBaseURL.plugins}/${pluginId}`, payload)
}

export function deletePlugin(pluginId: string) {
  return api.delete<null>(`${apiBaseURL.plugins}/${pluginId}`)
}

export type { PluginListQuery, PluginPayload, PluginRecord, UpdatePluginPayload } from './types'
