import { computed, h, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import {
  fetchProjectById,
  fetchProjectServices,
  fetchServiceDeployments,
  fetchServiceReleases,
  fetchTeams,
} from '@/api'
import { getEnvironmentTagColor } from '@/theme/environment'
import type { Deployment, Project, Release, Service, TeamRecord } from '@/api'
import { ApiError } from '@/api/request'

type ServiceReleaseRow = Release & { service_name: string }
type ServiceDeploymentRow = Deployment & { service_name: string }

export function useProjectDetailService() {
  const route = useRoute()
  const router = useRouter()
  const message = useMessage()

  const projectId = computed(() => route.params.projectId as string)
  const loading = ref(false)
  const project = ref<Project | null>(null)
  const teams = ref<TeamRecord[]>([])
  const services = ref<Service[]>([])
  const releases = ref<ServiceReleaseRow[]>([])
  const deployments = ref<ServiceDeploymentRow[]>([])

  const successfulReleases = computed(() =>
    releases.value.filter(item => item.status === 'completed').length,
  )

  const failedReleases = computed(() =>
    releases.value.filter(item => item.status === 'failed').length,
  )

  const successfulDeployments = computed(() =>
    deployments.value.filter(item => item.status === 'completed').length,
  )

  const failedDeployments = computed(() =>
    deployments.value.filter(item => item.status === 'failed').length,
  )

  const teamNameById = computed(() =>
    new Map(teams.value.map(team => [team.id, team.name])),
  )

  const teamOwnerContactById = computed(() =>
    new Map(teams.value.map(team => [team.id, team.owner_contact])),
  )

  const ownerTeamName = computed(() =>
    project.value?.team_id ? (teamNameById.value.get(project.value.team_id) || project.value.team_id) : 'Not set',
  )

  const ownerContact = computed(() =>
    project.value?.team_id ? (teamOwnerContactById.value.get(project.value.team_id) || 'Not set') : 'Not set',
  )

  const openProjects = () => {
    router.push({ name: 'projects' })
  }

  const openService = (row: Service) => {
    router.push({
      name: 'service-details',
      params: {
        serviceId: row.id,
      },
    })
  }

  const serviceColumns: DataTableColumns<Service> = [
    { title: 'Service', key: 'name' },
    {
      title: 'Repository',
      key: 'repo_url',
      render: row =>
        h(
          'a',
          {
            href: row.repo_url,
            target: '_blank',
            rel: 'noreferrer',
            class: 'text-[var(--app-accent)] hover:underline',
          },
          row.repo_url,
        ),
    },
  ]

  const releaseColumns: DataTableColumns<ServiceReleaseRow> = [
    { title: 'Service', key: 'service_name' },
    { title: 'Tag', key: 'tag' },
    { title: 'Target', key: 'target' },
    {
      title: 'Status',
      key: 'status',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: row.status === 'failed'
              ? { color: '#fee2e2', textColor: '#b91c1c' }
              : { color: '#dbeafe', textColor: '#1d4ed8' },
          },
          { default: () => row.status || 'pending' },
        ),
    },
  ]

  const deploymentColumns: DataTableColumns<ServiceDeploymentRow> = [
    { title: 'Service', key: 'service_name' },
    {
      title: 'Environment',
      key: 'environment',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: getEnvironmentTagColor(row.environment),
          },
          { default: () => row.environment },
        ),
    },
    { title: 'Version', key: 'version' },
    {
      title: 'Status',
      key: 'status',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: row.status === 'failed'
              ? { color: '#fee2e2', textColor: '#b91c1c' }
              : { color: '#dbeafe', textColor: '#1d4ed8' },
          },
          { default: () => row.status || 'pending' },
        ),
    },
  ]

  const loadProjectDetails = async() => {
    loading.value = true

    try {
      const [projectData, serviceData] = await Promise.all([
        fetchProjectById(projectId.value),
        fetchProjectServices(projectId.value),
      ])
      teams.value = await fetchTeams()

      project.value = projectData
      services.value = serviceData

      const serviceHistories = await Promise.all(
        serviceData.map(async (service) => {
          const [serviceReleases, serviceDeployments] = await Promise.all([
            fetchServiceReleases(service.id),
            fetchServiceDeployments(service.id, { limit: 5, sortBy: 'date', sortOrder: 'desc' }),
          ])

          return {
            releases: serviceReleases.map((item) => ({ ...item, service_name: service.name })),
            deployments: serviceDeployments.map((item) => ({ ...item, service_name: service.name })),
          }
        }),
      )

      releases.value = serviceHistories.flatMap(item => item.releases).slice(0, 8)
      deployments.value = serviceHistories.flatMap(item => item.deployments).slice(0, 8)
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load project details.')
    } finally {
      loading.value = false
    }
  }

  onMounted(loadProjectDetails)

  return {
    deploymentColumns,
    deployments,
    failedDeployments,
    failedReleases,
    getEnvironmentTagColor,
    loadProjectDetails,
    loading,
    openProjects,
    openService,
    ownerContact,
    ownerTeamName,
    project,
    releaseColumns,
    releases,
    serviceColumns,
    services,
    successfulDeployments,
    successfulReleases,
  }
}
