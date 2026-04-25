import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { permission } from '@/services/access/rbac'
import { fetchProjects, fetchTeams } from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import { environmentOptions, getEnvironmentTagColor } from '@/theme/environment'
import type { Project, TeamRecord } from '@/api'
import { ApiError } from '@/api/request'

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Active', value: 'active' },
  { label: 'Archived', value: 'archived' },
  { label: 'Deprecated', value: 'deprecated' },
]

export function useProjectListService() {
  const message = useMessage()
  const router = useRouter()
  const authStore = useAuthStore()
  const loading = ref(false)
  const rows = ref<Project[]>([])
  const teams = ref<TeamRecord[]>([])
  const filters = reactive({
    keyword: '',
    status: null as string | null,
    environment: null as string | null,
    ownerTeam: null as string | null,
  })

  const environmentSelectOptions = environmentOptions.map(option => ({ ...option }))

  const teamNameById = computed(() =>
    new Map(teams.value.map(team => [team.id, team.name])),
  )

  const getOwnerTeamName = (teamId?: string) => {
    if (!teamId) return ''
    return teamNameById.value.get(teamId) || teamId
  }

  const ownerTeamOptions = computed(() =>
    [...new Set(rows.value.map(row => getOwnerTeamName(row.team_id)))]
      .filter(Boolean)
      .map(value => ({ label: value, value })),
  )

  const canCreateProject = computed(() =>
    authStore.canAccess({ permissions: [permission.projectWrite] }),
  )

  const filteredRows = computed(() => {
    const keyword = filters.keyword.trim().toLowerCase()

    return rows.value.filter((row) => {
      const matchesKeyword = !keyword || [
        row.name,
        row.description,
        getOwnerTeamName(row.team_id),
        row.status,
      ].some(value => value?.toLowerCase().includes(keyword))

      const matchesStatus = !filters.status || row.status === filters.status
      const matchesEnvironment = !filters.environment || row.environments.includes(filters.environment)
      const matchesOwnerTeam = !filters.ownerTeam || getOwnerTeamName(row.team_id) === filters.ownerTeam

      return matchesKeyword && matchesStatus && matchesEnvironment && matchesOwnerTeam
    })
  })

  const openProject = (row: Project) => {
    router.push({ name: 'project-details', params: { projectId: row.id } })
  }

  const openProjectCreate = () => {
    router.push({ name: 'project-create' })
  }

  const resetFilters = () => {
    filters.keyword = ''
    filters.status = null
    filters.environment = null
    filters.ownerTeam = null
  }

  const columns: DataTableColumns<Project> = [
    { title: 'Name', key: 'name' },
    {
      title: 'Status',
      key: 'status',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: { color: '#dbeafe', textColor: '#1d4ed8' },
          },
          { default: () => row.status || 'Unknown' },
        ),
    },
    {
      title: 'Environments',
      key: 'environments',
      render: row =>
        h(
          'div',
          { class: 'flex flex-wrap gap-2' },
          row.environments.map((value) =>
            h(
              NTag,
              { bordered: false, color: getEnvironmentTagColor(value) },
              { default: () => value },
            ),
          ),
        ),
    },
    { title: 'Owner Team', key: 'owner_team', render: row => getOwnerTeamName(row.team_id) || 'Not set' },
    { title: 'Description', key: 'description' },
    {
      title: 'Actions',
      key: 'actions',
      render: row =>
        h(
          NButton,
          {
            size: 'small',
            secondary: false,
            onClick: (event: MouseEvent) => {
              event.stopPropagation()
              openProject(row)
            },
          },
          { default: () => 'View details' },
        ),
    },
  ]

  const loadProjects = async() => {
    loading.value = true
    try {
      const [projectData, teamData] = await Promise.all([
        fetchProjects(),
        fetchTeams(),
      ])
      rows.value = projectData
      teams.value = teamData
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load projects.')
    } finally {
      loading.value = false
    }
  }

  onMounted(loadProjects)

  return {
    canCreateProject,
    columns,
    environmentSelectOptions,
    filteredRows,
    filters,
    loadProjects,
    loading,
    openProject,
    openProjectCreate,
    ownerTeamOptions,
    resetFilters,
    statusOptions,
  }
}
