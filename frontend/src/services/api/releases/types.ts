export interface Release {
  id: string
  service_id: string
  plugin_id: string
  tag: string
  target: string
  name: string
  notes: string
  html_url: string
  external_ref: string
  status?: string
  triggered_by: string
}

export interface CreateReleasePayload {
  plugin_id: string
  tag: string
  target: string
  name?: string
  notes?: string
}
