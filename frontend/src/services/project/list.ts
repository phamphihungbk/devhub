import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { permission } from '@/services/access/rbac'
import { fetchProjects, fetchTeams } from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import type { Project, TeamRecord } from '@/api'
import { ApiError } from '@/api/request'

export function useProjectListService() {
  const message = useMessage()
  const router = useRouter()
  const authStore = useAuthStore()
  const loading = ref(false)
  const rows = ref<Project[]>([])
  const teams = ref<TeamRecord[]>([])
  const filters = reactive({
    keyword: '',
    ownerTeam: null as string | null,
  })

  const teamNameById = computed(() =>
    new Map(teams.value.map(team => [team.id, team.name])),
  )

  const getOwnerTeamName = (project: Project) => {
    const teamId = project.owner_team_id || project.team_id
    if (!teamId) return project.owner_team_name || project.owner_team || ''
    return project.owner_team_name || project.owner_team || teamNameById.value.get(teamId) || teamId
  }

  const shortId = (value?: string) => {
    if (!value) return 'Not set'
    return `${value.slice(0, 8)}...${value.slice(-4)}`
  }

  const ownerTeamOptions = computed(() =>
    [...new Set(rows.value.map(row => getOwnerTeamName(row)))]
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
        getOwnerTeamName(row),
        row.created_by_name,
        row.created_by,
      ].some(value => value?.toLowerCase().includes(keyword))

      const matchesOwnerTeam = !filters.ownerTeam || getOwnerTeamName(row) === filters.ownerTeam

      return matchesKeyword && matchesOwnerTeam
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
    filters.ownerTeam = null
  }

  const columns: DataTableColumns<Project> = [
    { title: 'Name', key: 'name' },
    { title: 'Owner Team', key: 'owner_team', render: row => getOwnerTeamName(row) || 'Not set' },
    {
      title: 'Created By',
      key: 'created_by',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: { color: '#f1f5f9', textColor: '#475569' },
          },
          { default: () => row.created_by_name || shortId(row.created_by) },
        ),
    },
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
    filteredRows,
    filters,
    loadProjects,
    loading,
    openProject,
    openProjectCreate,
    ownerTeamOptions,
    resetFilters,
  }
}
