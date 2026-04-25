import { computed, h, onMounted, ref } from 'vue'
import {useMessage } from 'naive-ui'


import { fetchPlugins, fetchProjects, fetchUsers } from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import { getRoleTagColor } from '@/theme/role'
import type { PluginRecord, Project, UserRecord } from '@/api'
import {useDashboardData} from './data'
import { ApiError } from '@/api/request'

export function useDashboardService() {
  const message = useMessage()
  const authStore = useAuthStore()
  const loading = ref(false)
  const projects = ref<Project[]>([])
  const plugins = ref<PluginRecord[]>([])
  const users = ref<UserRecord[]>([])
  const {pluginColumns, projectColumns} = useDashboardData()

  const teamMembers = computed(() =>
    users.value.filter(user => user.team_id === authStore.profile?.team_id),
  )

  const stats = computed(() => [
    {
      label: 'Projects',
      value: projects.value.length,
      caption: 'Tracked service domains',
    },
    {
      label: 'Plugins',
      value: plugins.value.length,
      caption: 'Installed platform automations',
    },
    {
      label: 'Users',
      value: teamMembers.value.length,
      caption: 'People in your team',
    },
  ])
  
  const loadDashboard = async() => {
    loading.value = true
    try {
      const [projectData, pluginData, userData] = await Promise.all([
        fetchProjects(),
        fetchPlugins(),
        fetchUsers(),
      ])
      projects.value = projectData
      plugins.value = pluginData
      users.value = userData
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load dashboard data.')
    } finally {
      loading.value = false
    }
  }

  onMounted(loadDashboard)

  return {
    getRoleTagColor,
    loadDashboard,
    loading,
    pluginColumns,
    plugins,
    projectColumns,
    projects,
    stats,
    teamMembers,
  }
}
