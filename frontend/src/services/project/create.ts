import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'

import { createProject } from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import { ApiError } from '@/api/request'

export function useProjectCreateService() {
  const router = useRouter()
  const message = useMessage()
  const authStore = useAuthStore()
  const saving = ref(false)
  const form = reactive({
    name: '',
    description: '',
  })

  const validateForm = () => {
    if (!form.name.trim()) return 'Project name is required.'
    if (!authStore.profile?.team_id) return 'Your team information is not available.'
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
        name: form.name.trim(),
        description: form.description?.trim() || undefined,
        owner_team_id: teamId,
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
    form,
    openProjects,
    saving,
    submitProject,
  }
}
