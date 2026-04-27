import { computed, reactive, ref, type Ref } from 'vue'
import { useMessage } from 'naive-ui'

import { permission } from '@/services/access/rbac'
import {
  createScaffoldRequest,
  suggestProjectScaffoldRequest,
} from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import type {
  CreateScaffoldRequestPayload,
  PluginRecord,
  Project,
  ScaffoldRequestSuggestion,
} from '@/api'
import { ApiError } from '@/api/request'

interface UseCreateScaffoldRequestServiceInput {
  projects: Ref<Project[]>
  plugins: Ref<PluginRecord[]>
  onCreated: () => Promise<void>
}

export interface GenerateScaffoldSuggestionInput {
  projectId: string
  prompt: string
}

export async function generateScaffoldSuggestion(input: GenerateScaffoldSuggestionInput) {
  return suggestProjectScaffoldRequest(input.projectId, {
    prompt: input.prompt.trim(),
  })
}

export function applyScaffoldSuggestionToForm(
  form: CreateScaffoldRequestPayload,
  suggestion: ScaffoldRequestSuggestion,
) {
  form.plugin_id = suggestion.plugin_id || form.plugin_id
  form.environment = suggestion.environment
  form.variables = { ...suggestion.variables }
}

export function useCreateScaffoldRequestService(input: UseCreateScaffoldRequestServiceInput) {
  const message = useMessage()
  const authStore = useAuthStore()
  const submitting = ref(false)
  const suggestionLoading = ref(false)
  const modalOpen = ref(false)
  const suggestion = ref<ScaffoldRequestSuggestion | null>(null)
  const prompt = ref('')
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
    input.projects.value.map(project => ({ label: project.name, value: project.id })),
  )

  const selectedProject = computed(() =>
    input.projects.value.find(project => project.id === form.project_id) || null,
  )

  const environmentOptions = computed(() => {
    const values = selectedProject.value?.environments?.length
      ? selectedProject.value.environments
      : ['dev', 'staging', 'prod']
    return values.map(value => ({ label: value, value }))
  })

  const scaffolderOptions = computed(() =>
    input.plugins.value
      .filter(plugin => plugin.type === 'scaffolder' && plugin.enabled !== false)
      .map(plugin => ({ label: plugin.name, value: plugin.id })),
  )

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

    applyScaffoldSuggestionToForm(form, suggestion.value)
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
      suggestion.value = await generateScaffoldSuggestion({
        projectId: form.project_id,
        prompt: prompt.value,
      })
      applySuggestion()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to suggest scaffold request.')
    } finally {
      suggestionLoading.value = false
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
      await input.onCreated()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to create scaffold request.')
    } finally {
      submitting.value = false
    }
  }

  return {
    analyzePrompt,
    applySuggestion,
    canCreateScaffoldRequest,
    environmentOptions,
    form,
    handleProjectChange,
    modalOpen,
    openModal,
    projectOptions,
    prompt,
    scaffolderOptions,
    submitting,
    submitScaffoldRequest,
    suggestion,
    suggestionLoading,
  }
}
