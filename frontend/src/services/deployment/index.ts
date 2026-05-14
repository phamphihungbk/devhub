import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { permission } from '@/services/access/rbac'
import { createDeployment, fetchDeploymentById, fetchPlugins, fetchProjectById, fetchProjects, fetchProjectServices, fetchServiceDeployments, fetchServiceReleases } from '@/api'
import { getEnvironmentTagColor } from '@/theme/environment'
import { useAuthStore } from '@/stores/modules/auth'
import type { CreateDeploymentPayload, Deployment, PluginRecord, Project, ProjectEnvironment, Release, Service } from '@/api'
import { ApiError } from '@/api/request'

export type DeploymentRow = Deployment & {
  project_id: string
  project_name: string
  service_name: string
}

type ServiceOptionRecord = Service & {
  project_name: string
  environments: ProjectEnvironment[]
}

function projectEnvironments(project: Project): ProjectEnvironment[] {
  return (project.environments || [])
    .filter((environment): environment is ProjectEnvironment => typeof environment !== 'string')
}

export function getDeploymentStatusTagColor(status?: string) {
  switch (status) {
    case 'completed':
      return { color: '#dcfce7', textColor: '#15803d' }
    case 'failed':
      return { color: '#fee2e2', textColor: '#b91c1c' }
    case 'running':
      return { color: '#fef3c7', textColor: '#b45309' }
    default:
      return { color: '#dbeafe', textColor: '#1d4ed8' }
  }
}

export function formatRunnerText(value?: string) {
  return value?.trim() || 'No runner output recorded yet.'
}

export function useDeploymentListService() {
  const message = useMessage()
  const router = useRouter()
  const authStore = useAuthStore()
  const loading = ref(false)
  const logLoading = ref(false)
  const deploymentSubmitting = ref(false)
  const rows = ref<DeploymentRow[]>([])
  const services = ref<ServiceOptionRecord[]>([])
  const plugins = ref<PluginRecord[]>([])
  const releases = ref<Release[]>([])
  const deploymentReleases = ref<Release[]>([])
  const selectedDeployment = ref<DeploymentRow | null>(null)
  const deploymentModalOpen = ref(false)
  const logModalOpen = ref(false)
  const filters = reactive({
    keyword: '',
    environment: null as string | null,
    status: null as string | null,
  })
  const deploymentForm = reactive<CreateDeploymentPayload & { service_id: string }>({
    service_id: '',
    plugin_id: '',
    environment_id: '',
    release_id: '',
  })

  const canCreateDeployment = computed(() =>
    authStore.canAccess({ permissions: [permission.deploymentWrite] }),
  )

  const environmentOptions = computed(() =>
    [...new Set(rows.value.map(row => deploymentEnvironmentName(row)).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const statusOptions = computed(() =>
    [...new Set(rows.value.map(row => row.status).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const serviceOptions = computed(() =>
    services.value.map(service => ({
      label: `${service.name} · ${service.project_name}`,
      value: service.id,
    })),
  )

  const selectedService = computed(() =>
    services.value.find(service => service.id === deploymentForm.service_id) || null,
  )

  const deployerOptions = computed(() =>
    plugins.value
      .filter(plugin => plugin.type === 'deployer' && plugin.enabled !== false)
      .map(plugin => ({ label: plugin.name, value: plugin.id })),
  )

  const deploymentEnvironmentOptions = computed(() =>
    (selectedService.value?.environments || []).map(environment => ({
      label: environment.name,
      value: environment.id,
    })),
  )

  const releaseOptions = computed(() =>
    deploymentReleases.value
      .filter(release => release.status === 'completed' || !release.status)
      .map(release => ({
        label: `${release.tag}${release.name ? ` · ${release.name}` : ''}`,
        value: release.id,
      })),
  )

  const environmentNameById = computed(() => {
    const entries = services.value.flatMap(service =>
      service.environments.map(environment => [environment.id, environment.name] as const),
    )
    return new Map(entries)
  })

  const releaseTagById = computed(() =>
    new Map(releases.value.map(release => [release.id, release.tag])),
  )

  const deploymentEnvironmentName = (row: Deployment) =>
    row.environment || (row.environment_id ? environmentNameById.value.get(row.environment_id) : '') || 'Unknown'

  const deploymentVersion = (row: Deployment) =>
    row.version || (row.release_id ? releaseTagById.value.get(row.release_id) : '') || 'Unknown'

  const filteredRows = computed(() => {
    const keyword = filters.keyword.trim().toLowerCase()

    return rows.value.filter((row) => {
      const matchesKeyword = !keyword || [
        row.project_name,
        row.service_name,
        deploymentEnvironmentName(row),
        deploymentVersion(row),
        row.status,
        row.commit_sha,
        row.external_ref,
      ].some(value => value?.toLowerCase().includes(keyword))
      const matchesEnvironment = !filters.environment || deploymentEnvironmentName(row) === filters.environment
      const matchesStatus = !filters.status || row.status === filters.status

      return matchesKeyword && matchesEnvironment && matchesStatus
    })
  })

  const runningCount = computed(() =>
    rows.value.filter(row => row.status === 'running').length,
  )

  const completedCount = computed(() =>
    rows.value.filter(row => row.status === 'completed').length,
  )

  const failedCount = computed(() =>
    rows.value.filter(row => row.status === 'failed').length,
  )

  const openService = (row: DeploymentRow) => {
    router.push({
      name: 'service-details',
      params: {
        serviceId: row.service_id,
      },
    })
  }

  const openLogs = async(row: DeploymentRow) => {
    selectedDeployment.value = row
    logModalOpen.value = true
    logLoading.value = true

    try {
      const deployment = await fetchDeploymentById(row.id)
      selectedDeployment.value = {
        ...row,
        ...deployment,
      }
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load deployment output.')
    } finally {
      logLoading.value = false
    }
  }

  const resetFilters = () => {
    filters.keyword = ''
    filters.environment = null
    filters.status = null
  }

  const resetDeploymentForm = () => {
    deploymentForm.service_id = serviceOptions.value[0]?.value || ''
    deploymentForm.plugin_id = deployerOptions.value[0]?.value || ''
    deploymentForm.environment_id = deploymentEnvironmentOptions.value[0]?.value || ''
    deploymentForm.release_id = releaseOptions.value[0]?.value || ''
  }

  const openDeploymentModal = async() => {
    resetDeploymentForm()
    deploymentModalOpen.value = true
    await loadSelectedServiceReleases()
  }

  const loadSelectedServiceReleases = async() => {
    deploymentReleases.value = deploymentForm.service_id
      ? await fetchServiceReleases(deploymentForm.service_id)
      : []
    deploymentForm.release_id = releaseOptions.value[0]?.value || ''
  }

  const handleDeploymentServiceChange = async() => {
    deploymentForm.environment_id = deploymentEnvironmentOptions.value[0]?.value || ''
    await loadSelectedServiceReleases()
  }

  const columns: DataTableColumns<DeploymentRow> = [
    { title: 'Service', key: 'service_name' },
    { title: 'Project', key: 'project_name' },
    {
      title: 'Environment',
      key: 'environment',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: getEnvironmentTagColor(deploymentEnvironmentName(row)),
          },
          { default: () => deploymentEnvironmentName(row) },
        ),
    },
    { title: 'Version', key: 'version', render: row => deploymentVersion(row) },
    {
      title: 'Status',
      key: 'status',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: getDeploymentStatusTagColor(row.status),
          },
          { default: () => row.status || 'pending' },
        ),
    },
    { title: 'Commit SHA', key: 'commit_sha', render: row => row.commit_sha || 'Not set' },
    { title: 'Finished', key: 'finished_at', render: row => formatOptionalDate(row.finished_at) },
    {
      title: 'Actions',
      key: 'actions',
      render: row =>
        h(
          'div',
          { class: 'flex gap-2' },
          [
            h(
              NButton,
              {
                size: 'small',
                onClick: (event: MouseEvent) => {
                  event.stopPropagation()
                  openLogs(row)
                },
              },
              { default: () => 'Output' },
            ),
            h(
              NButton,
              {
                size: 'small',
                onClick: (event: MouseEvent) => {
                  event.stopPropagation()
                  openService(row)
                },
              },
              { default: () => 'Service' },
            ),
          ],
        ),
    },
  ]

  const loadDeployments = async() => {
    loading.value = true
    try {
      releases.value = []
      const [projects, pluginRows] = await Promise.all([
        fetchProjects(),
        fetchPlugins(),
      ])
      const deploymentGroups = await Promise.all(
        projects.map(async (project: Project) => {
          const [projectDetail, projectServices] = await Promise.all([
            fetchProjectById(project.id),
            fetchProjectServices(project.id),
          ])
          const serviceRows = projectServices.map(service => ({
            ...service,
            project_name: project.name,
            environments: projectEnvironments(projectDetail),
          }))

          const deploymentRows = await Promise.all(
            projectServices.map(async (service: Service) => {
              const [deployments, serviceReleases] = await Promise.all([
                fetchServiceDeployments(service.id, { limit: 50, sortBy: 'date', sortOrder: 'desc' }),
                fetchServiceReleases(service.id),
              ])
              releases.value.push(...serviceReleases)

              return deployments.map((deployment: Deployment) => ({
                ...deployment,
                project_id: project.id,
                project_name: project.name,
                service_name: service.name,
              }))
            }),
          )

          return {
            services: serviceRows,
            deployments: deploymentRows.flat(),
          }
        }),
      )

      services.value = deploymentGroups.flatMap(group => group.services)
      plugins.value = pluginRows
      rows.value = deploymentGroups.flatMap(group => group.deployments)
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load deployments.')
    } finally {
      loading.value = false
    }
  }

  const submitDeployment = async() => {
    if (!deploymentForm.service_id || !deploymentForm.plugin_id || !deploymentForm.environment_id) {
      message.warning('Complete the deployment form before submitting.')
      return
    }

    deploymentSubmitting.value = true
    try {
      await createDeployment(deploymentForm.service_id, {
        plugin_id: deploymentForm.plugin_id,
        environment_id: deploymentForm.environment_id,
        release_id: deploymentForm.release_id || undefined,
      })
      message.success('Deployment created successfully.')
      deploymentModalOpen.value = false
      await loadDeployments()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to create deployment.')
    } finally {
      deploymentSubmitting.value = false
    }
  }

  onMounted(loadDeployments)

  return {
    canCreateDeployment,
    columns,
    completedCount,
    deployerOptions,
    deploymentEnvironmentOptions,
    deploymentForm,
    deploymentModalOpen,
    deploymentSubmitting,
    environmentOptions,
    failedCount,
    filteredRows,
    filters,
    formatRunnerText,
    deploymentEnvironmentName,
    deploymentVersion,
    getDeploymentStatusTagColor,
    handleDeploymentServiceChange,
    loadDeployments,
    loading,
    logLoading,
    logModalOpen,
    openDeploymentModal,
    openLogs,
    resetFilters,
    rows,
    runningCount,
    selectedDeployment,
    serviceOptions,
    statusOptions,
    releaseOptions,
    submitDeployment,
  }
}
