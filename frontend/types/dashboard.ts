export type HighlightTone = 'accent' | 'positive' | 'neutral'
export type StatusTone = 'accent' | 'positive' | 'neutral' | 'warning'

export interface DashboardProject {
  id: string
  name: string
  description: string
  environments: string[]
  createdBy: string
}

export interface DashboardPlugin {
  id: string
  name: string
  type: string
  version: string
  entrypoint: string
  scope: string
  description: string
}

export interface DashboardUser {
  id: string
  name: string
  email: string
  role: string
}

export interface DashboardDeployment {
  id: string
  projectId: string
  environment: string
  service: string
  version: string
  status: string
  statusTone: StatusTone
  triggeredBy: string
}

export interface DashboardHighlight {
  label: string
  value: string
  description: string
  tone: HighlightTone
}

export interface DashboardOverview {
  generatedAt: string
  highlights: DashboardHighlight[]
  projects: DashboardProject[]
  plugins: DashboardPlugin[]
  users: DashboardUser[]
  deployments: DashboardDeployment[]
}
