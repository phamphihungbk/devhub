import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { fetchProjects, fetchProjectServices } from '@/api'
import type { Project, Service } from '@/api'
import { ApiError } from '@/api/request'

export type ServiceRow = Service & {
  project_name: string
}

export function useServiceListService() {
  const message = useMessage()
  const router = useRouter()
  const loading = ref(false)
  const rows = ref<ServiceRow[]>([])
  const filters = reactive({
    keyword: '',
  })

  const filteredRows = computed(() => {
    const keyword = filters.keyword.trim().toLowerCase()

    return rows.value.filter((row) => {
      return !keyword || [
        row.name,
        row.repo_url,
        row.project_name,
      ].some(value => value?.toLowerCase().includes(keyword))
    })
  })

  const openService = (row: ServiceRow) => {
    router.push({
      name: 'service-details',
      params: {
        serviceId: row.id,
      },
    })
  }

  const resetFilters = () => {
    filters.keyword = ''
  }

  const columns: DataTableColumns<ServiceRow> = [
    { title: 'Service', key: 'name' },
    { title: 'Project', key: 'project_name' },
    {
      title: 'Repository',
      key: 'repo_url',
      render: row =>
        h(
          'a',
          {
            href: row.repo_url,
            target: '_blank',
            rel: 'noreferrer',
            class: 'text-[var(--app-accent)] hover:underline',
          },
          row.repo_url,
        ),
    },
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
              openService(row)
            },
          },
          { default: () => 'View details' },
        ),
    },
  ]

  const loadServices = async() => {
    loading.value = true
    try {
      const projects = await fetchProjects()
      const serviceGroups = await Promise.all(
        projects.map(async (project: Project) => {
          const services = await fetchProjectServices(project.id)
          return services.map((service) => ({
            ...service,
            project_name: project.name,
          }))
        }),
      )

      rows.value = serviceGroups.flat()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load services.')
    } finally {
      loading.value = false
    }
  }

  onMounted(loadServices)

  return {
    columns,
    filteredRows,
    filters,
    loadServices,
    loading,
    openService,
    resetFilters,
  }
}
