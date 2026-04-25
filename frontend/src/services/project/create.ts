import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'

import { createProject } from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import { environmentOptions } from '@/theme/environment'
import type { ProjectPayload } from '@/api'
import { ApiError } from '@/api/request'

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Active', value: 'active' },
  { label: 'Archived', value: 'archived' },
  { label: 'Deprecated', value: 'deprecated' },
]

const scmProviderOptions = [
  { label: 'Gitea', value: 'gitea' },
  { label: 'GitHub', value: 'github' },
  { label: 'GitLab', value: 'gitlab' },
  { label: 'Bitbucket', value: 'bitbucket' },
]

export function useProjectCreateService() {
  const router = useRouter()
  const message = useMessage()
  const authStore = useAuthStore()
  const saving = ref(false)
  const environmentSelectOptions = environmentOptions.map(option => ({ ...option }))
  const form = reactive<ProjectPayload>({
    name: '',
    description: '',
    environments: ['dev'],
    status: 'draft',
    team_id: '',
    scm_provider: 'gitea',
  })

  const validateForm = () => {
    if (!form.name.trim()) return 'Project name is required.'
    if (form.environments.length === 0) return 'Select at least one environment.'
    if (!authStore.profile?.team_id) return 'Your team information is not available.'
    if (!form.scm_provider.trim()) return 'SCM provider is required.'
    return null
  }

  const submitProject = async() => {
    const validationError = validateForm()
    const teamId = authStore.profile?.team_id

    if (validationError || !teamId) {
      message.warning(validationError || 'Your team information is not available.')
      return
    }

    saving.value = true

    try {
      await createProject({
        ...form,
        team_id: teamId,
        name: form.name.trim(),
        description: form.description?.trim() || undefined,
        scm_provider: form.scm_provider.trim(),
      })

      message.success('Project created successfully.')
      await router.push({ name: 'projects' })
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to create project.')
    } finally {
      saving.value = false
    }
  }

  const openProjects = () => {
    router.push({ name: 'projects' })
  }

  return {
    environmentSelectOptions,
    form,
    openProjects,
    saving,
    scmProviderOptions,
    statusOptions,
    submitProject,
  }
}
