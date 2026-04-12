import type {
  DashboardDeployment,
  DashboardHighlight,
  DashboardOverview,
  DashboardPlugin,
  DashboardProject,
  DashboardUser,
  StatusTone,
} from '~/types/dashboard'
import { fetchBackend } from '../../utils/backend'

interface BackendProject {
  id: string
  name: string
  description: string
  environments: string[]
  created_by: string
}

interface BackendPlugin {
  id: string
  name: string
  type: string
  version: string
  entrypoint: string
  scope: string
  description: string
}

interface BackendUser {
  id: string
  name: string
  email: string
  role: string
}

interface BackendDeployment {
  id: string
  project_id: string
  environment: string
  service: string
  version: string
  status: string
  triggered_by: string
}

const statusToneMap: Record<string, StatusTone> = {
  success: 'positive',
  completed: 'positive',
  running: 'accent',
  in_progress: 'accent',
  pending: 'neutral',
  queued: 'neutral',
  failed: 'warning',
  error: 'warning',
}

function toProject(project: BackendProject): DashboardProject {
  return {
    id: project.id,
    name: project.name,
    description: project.description,
    environments: project.environments,
    createdBy: project.created_by,
  }
}

function toPlugin(plugin: BackendPlugin): DashboardPlugin {
  return {
    id: plugin.id,
    name: plugin.name,
    type: plugin.type,
    version: plugin.version,
    entrypoint: plugin.entrypoint,
    scope: plugin.scope,
    description: plugin.description,
  }
}

function toUser(user: BackendUser): DashboardUser {
  return {
    id: user.id,
    name: user.name,
    email: user.email,
    role: user.role,
  }
}

function toDeployment(deployment: BackendDeployment): DashboardDeployment {
  const normalizedStatus = deployment.status.toLowerCase()

  return {
    id: deployment.id,
    projectId: deployment.project_id,
    environment: deployment.environment,
    service: deployment.service,
    version: deployment.version,
    status: deployment.status,
    statusTone: statusToneMap[normalizedStatus] || 'neutral',
    triggeredBy: deployment.triggered_by,
  }
}

export default defineEventHandler(async (): Promise<DashboardOverview> => {
  const [projectData, pluginData, userData] = await Promise.all([
    fetchBackend<BackendProject[]>('/projects/'),
    fetchBackend<BackendPlugin[]>('/plugins/'),
    fetchBackend<BackendUser[]>('/users/'),
  ])

  const projects = projectData.map(toProject)
  const plugins = pluginData.map(toPlugin)
  const users = userData.map(toUser)

  const deploymentsByProject = await Promise.all(
    projects.slice(0, 3).map(async (project) => {
      try {
        const data = await fetchBackend<BackendDeployment[]>(`/projects/${project.id}/deployments`)
        return data.map(toDeployment)
      } catch {
        return []
      }
    }),
  )

  const deployments = deploymentsByProject.flat().slice(0, 6)
  const highlights: DashboardHighlight[] = [
    {
      label: 'Projects',
      value: String(projects.length),
      description: 'Active service spaces connected to the platform.',
      tone: 'accent',
    },
    {
      label: 'Plugins',
      value: String(plugins.length),
      description: 'Installed extensions shaping scaffolding and workflows.',
      tone: 'positive',
    },
    {
      label: 'Users',
      value: String(users.length),
      description: 'Known teammates with access to the control plane.',
      tone: 'neutral',
    },
    {
      label: 'Recent Deployments',
      value: String(deployments.length),
      description: 'Latest deployment records pulled from the backend.',
      tone: deployments.length > 0 ? 'positive' : 'neutral',
    },
  ]

  return {
    generatedAt: new Date().toISOString(),
    highlights,
    projects,
    plugins,
    users,
    deployments,
  }
})
