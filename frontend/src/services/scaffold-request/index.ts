import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { permission } from '@/services/access/rbac'
import {
  createScaffoldRequest,
  fetchPlugins,
  fetchProjectScaffoldRequests,
  fetchProjects,
  suggestProjectScaffoldRequest,
} from '@/api'
import { getEnvironmentTagColor } from '@/theme/environment'
import { useAuthStore } from '@/stores/modules/auth'
import type {
  CreateScaffoldRequestPayload,
  PluginRecord,
  Project,
  ScaffoldRequestRecord,
  ScaffoldRequestSuggestion,
} from '@/api'
import { ApiError } from '@/api/request'

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
  const authStore = useAuthStore()
  const loading = ref(false)
  const submitting = ref(false)
  const suggestionLoading = ref(false)
  const modalOpen = ref(false)
  const rows = ref<ScaffoldRequestRow[]>([])
  const projects = ref<Project[]>([])
  const plugins = ref<PluginRecord[]>([])
  const suggestion = ref<ScaffoldRequestSuggestion | null>(null)
  const prompt = ref('')
  const filters = reactive({
    keyword: '',
    status: null as string | null,
    environment: null as string | null,
  })
  const form = reactive<CreateScaffoldRequestPayload & { project_id: string }>({
    project_id: '',
    plugin_id: '',
    environment: 'dev',
    variables: {
      service_name: '',
      module_path: '',
      port: 8080,
      database: 'postgres',
      enable_logging: true,
    },
  })

  const canCreateScaffoldRequest = computed(() =>
    authStore.canAccess({ permissions: [permission.scaffoldRequestWrite] }),
  )

  const projectOptions = computed(() =>
    projects.value.map(project => ({ label: project.name, value: project.id })),
  )

  const selectedProject = computed(() =>
    projects.value.find(project => project.id === form.project_id) || null,
  )

  const environmentOptions = computed(() => {
    const values = selectedProject.value?.environments?.length
      ? selectedProject.value.environments
      : ['dev', 'staging', 'prod']
    return values.map(value => ({ label: value, value }))
  })

  const statusOptions = computed(() =>
    [...new Set(rows.value.map(row => row.status).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const filterEnvironmentOptions = computed(() =>
    [...new Set(rows.value.map(row => row.environment || row.environments).filter(Boolean))]
      .map(value => ({ label: value, value })),
  )

  const scaffolderOptions = computed(() =>
    plugins.value
      .filter(plugin => plugin.type === 'scaffolder')
      .map(plugin => ({ label: plugin.name, value: plugin.id })),
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

  const resetForm = () => {
    form.project_id = projectOptions.value[0]?.value || ''
    form.plugin_id = scaffolderOptions.value[0]?.value || ''
    form.environment = selectedProject.value?.environments?.[0] || 'dev'
    form.variables.service_name = ''
    form.variables.module_path = ''
    form.variables.port = 8080
    form.variables.database = 'postgres'
    form.variables.enable_logging = true
    prompt.value = ''
    suggestion.value = null
  }

  const openModal = () => {
    resetForm()
    modalOpen.value = true
  }

  const handleProjectChange = () => {
    form.environment = selectedProject.value?.environments?.[0] || 'dev'
    suggestion.value = null
  }

  const applySuggestion = () => {
    if (!suggestion.value) return

    form.plugin_id = suggestion.value.plugin_id || form.plugin_id
    form.environment = suggestion.value.environment
    form.variables = { ...suggestion.value.variables }
  }

  const analyzePrompt = async() => {
    if (!form.project_id) {
      message.warning('Select a project before analyzing the prompt.')
      return
    }
    if (!prompt.value.trim()) {
      message.warning('Describe the service before asking for a suggestion.')
      return
    }

    suggestionLoading.value = true
    try {
      suggestion.value = await suggestProjectScaffoldRequest(form.project_id, {
        prompt: prompt.value.trim(),
        project_name: selectedProject.value?.name || '',
        project_description: selectedProject.value?.description || '',
        environment: form.environment,
        environments: selectedProject.value?.environments || [],
      })
      applySuggestion()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to suggest scaffold request.')
    } finally {
      suggestionLoading.value = false
    }
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

  const submitScaffoldRequest = async() => {
    if (!form.project_id || !form.plugin_id || !form.environment || !form.variables.service_name.trim()) {
      message.warning('Complete project, plugin, environment, and service name before creating the request.')
      return
    }

    submitting.value = true
    try {
      await createScaffoldRequest(form.project_id, {
        plugin_id: form.plugin_id,
        environment: form.environment,
        variables: {
          ...form.variables,
          service_name: form.variables.service_name.trim(),
          module_path: form.variables.module_path.trim(),
        },
      })
      message.success('Scaffold request created successfully.')
      modalOpen.value = false
      await loadScaffoldRequests()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to create scaffold request.')
    } finally {
      submitting.value = false
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

  onMounted(loadScaffoldRequests)

  return {
    analyzePrompt,
    applySuggestion,
    canCreateScaffoldRequest,
    columns,
    completedCount,
    environmentOptions,
    failedCount,
    filterEnvironmentOptions,
    filteredRows,
    filters,
    form,
    handleProjectChange,
    loading,
    modalOpen,
    openModal,
    openProject,
    pendingCount,
    projectOptions,
    prompt,
    resetFilters,
    rows,
    runningCount,
    scaffolderOptions,
    statusOptions,
    submitting,
    submitScaffoldRequest,
    suggestion,
    suggestionLoading,
  }
}
