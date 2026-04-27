import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import {
  fetchPlugins,
  fetchProjectScaffoldRequests,
  fetchProjects,
} from '@/api'
import { getEnvironmentTagColor } from '@/theme/environment'
import type {
  PluginRecord,
  Project,
  ScaffoldRequestRecord,
} from '@/api'
import { ApiError } from '@/api/request'
import { useCreateScaffoldRequestService } from './create'

export type ScaffoldRequestRow = ScaffoldRequestRecord & {
  project_name: string
}

function getStatusTagColor(status: string) {
  switch (status) {
    case 'completed':
    case 'approved':
      return { color: '#dcfce7', textColor: '#15803d' }
    case 'failed':
    case 'rejected':
      return { color: '#fee2e2', textColor: '#b91c1c' }
    case 'running':
      return { color: '#dbeafe', textColor: '#1d4ed8' }
    default:
      return { color: '#fef3c7', textColor: '#b45309' }
  }
}

export function useScaffoldRequestListService() {
  const message = useMessage()
  const router = useRouter()
  const loading = ref(false)
  const rows = ref<ScaffoldRequestRow[]>([])
  const projects = ref<Project[]>([])
  const plugins = ref<PluginRecord[]>([])
  const filters = reactive({
    keyword: '',
    status: null as string | null,
    environment: null as string | null,
  })

  const statusOptions = computed(() =>
    [...new Set(rows.value.map(row => row.status).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const filterEnvironmentOptions = computed(() =>
    [...new Set(rows.value.map(row => row.environment || row.environments).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const pendingCount = computed(() =>
    rows.value.filter(row => row.status === 'pending').length,
  )

  const runningCount = computed(() =>
    rows.value.filter(row => row.status === 'running').length,
  )

  const completedCount = computed(() =>
    rows.value.filter(row => row.status === 'completed').length,
  )

  const failedCount = computed(() =>
    rows.value.filter(row => ['failed', 'rejected'].includes(row.status)).length,
  )

  const filteredRows = computed(() => {
    const keyword = filters.keyword.trim().toLowerCase()

    return rows.value.filter((row) => {
      const environment = row.environment || row.environments || ''
      const matchesKeyword = !keyword || [
        row.project_name,
        row.status,
        environment,
        row.variables?.service_name,
        row.variables?.module_path,
        row.variables?.database,
      ].some(value => value?.toLowerCase().includes(keyword))
      const matchesStatus = !filters.status || row.status === filters.status
      const matchesEnvironment = !filters.environment || environment === filters.environment

      return matchesKeyword && matchesStatus && matchesEnvironment
    })
  })

  const resetFilters = () => {
    filters.keyword = ''
    filters.status = null
    filters.environment = null
  }

  const loadScaffoldRequests = async() => {
    loading.value = true
    try {
      const [projectRows, pluginRows] = await Promise.all([
        fetchProjects(),
        fetchPlugins(),
      ])

      const requestGroups = await Promise.all(
        projectRows.map(async (project: Project) => {
          const requests = await fetchProjectScaffoldRequests(project.id, {
            limit: 50,
            sortBy: 'date',
            sortOrder: 'desc',
          })

          return requests.map((request: ScaffoldRequestRecord) => ({
            ...request,
            project_name: project.name,
          }))
        }),
      )

      projects.value = projectRows
      plugins.value = pluginRows
      rows.value = requestGroups.flat()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load scaffold requests.')
    } finally {
      loading.value = false
    }
  }

  const openProject = (row: ScaffoldRequestRow) => {
    router.push({ name: 'project-details', params: { projectId: row.project_id } })
  }

  const columns: DataTableColumns<ScaffoldRequestRow> = [
    { title: 'Service', key: 'service', render: row => row.variables?.service_name || 'Not set' },
    { title: 'Project', key: 'project_name' },
    {
      title: 'Environment',
      key: 'environment',
      render: (row) => {
        const environment = row.environment || row.environments || 'dev'
        return h(
          NTag,
          {
            bordered: false,
            color: getEnvironmentTagColor(environment),
          },
          { default: () => environment },
        )
      },
    },
    {
      title: 'Status',
      key: 'status',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: getStatusTagColor(row.status),
          },
          { default: () => row.status || 'pending' },
        ),
    },
    { title: 'Database', key: 'database', render: row => row.variables?.database || 'Not set' },
    { title: 'Port', key: 'port', render: row => row.variables?.port || 'Not set' },
    {
      title: 'Actions',
      key: 'actions',
      render: row =>
        h(
          NButton,
          {
            size: 'small',
            onClick: (event: MouseEvent) => {
              event.stopPropagation()
              openProject(row)
            },
          },
          { default: () => 'Project' },
        ),
    },
  ]

  const createService = useCreateScaffoldRequestService({
    projects,
    plugins,
    onCreated: loadScaffoldRequests,
  })

  onMounted(loadScaffoldRequests)

  return {
    columns,
    completedCount,
    failedCount,
    filterEnvironmentOptions,
    filteredRows,
    filters,
    loading,
    openProject,
    pendingCount,
    resetFilters,
    rows,
    runningCount,
    statusOptions,
    ...createService,
  }
}
