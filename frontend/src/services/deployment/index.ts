import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { fetchDeploymentById, fetchProjects, fetchProjectServices, fetchServiceDeployments } from '@/api'
import { getEnvironmentTagColor } from '@/theme/environment'
import type { Deployment, Project, Service } from '@/api'
import { ApiError } from '@/api/request'

export type DeploymentRow = Deployment & {
  project_id: string
  project_name: string
  service_name: string
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
  const loading = ref(false)
  const logLoading = ref(false)
  const rows = ref<DeploymentRow[]>([])
  const selectedDeployment = ref<DeploymentRow | null>(null)
  const logModalOpen = ref(false)
  const filters = reactive({
    keyword: '',
    environment: null as string | null,
    status: null as string | null,
  })

  const environmentOptions = computed(() =>
    [...new Set(rows.value.map(row => row.environment).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const statusOptions = computed(() =>
    [...new Set(rows.value.map(row => row.status).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const filteredRows = computed(() => {
    const keyword = filters.keyword.trim().toLowerCase()

    return rows.value.filter((row) => {
      const matchesKeyword = !keyword || [
        row.project_name,
        row.service_name,
        row.environment,
        row.version,
        row.status,
        row.commit_sha,
        row.external_ref,
      ].some(value => value?.toLowerCase().includes(keyword))
      const matchesEnvironment = !filters.environment || row.environment === filters.environment
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
      const projects = await fetchProjects()
      const deploymentGroups = await Promise.all(
        projects.map(async (project: Project) => {
          const services = await fetchProjectServices(project.id)

          return Promise.all(
            services.map(async (service: Service) => {
              const deployments = await fetchServiceDeployments(service.id, { limit: 50, sortBy: 'date', sortOrder: 'desc' })

              return deployments.map((deployment: Deployment) => ({
                ...deployment,
                project_id: project.id,
                project_name: project.name,
                service_name: service.name,
              }))
            }),
          )
        }),
      )

      rows.value = deploymentGroups.flatMap(group => group.flat())
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load deployments.')
    } finally {
      loading.value = false
    }
  }

  onMounted(loadDeployments)

  return {
    columns,
    completedCount,
    environmentOptions,
    failedCount,
    filteredRows,
    filters,
    formatRunnerText,
    getDeploymentStatusTagColor,
    loadDeployments,
    loading,
    logLoading,
    logModalOpen,
    openLogs,
    resetFilters,
    rows,
    runningCount,
    selectedDeployment,
    statusOptions,
  }
}
