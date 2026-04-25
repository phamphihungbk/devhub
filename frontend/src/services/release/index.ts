import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { permission } from '@/services/access/rbac'
import { createRelease, fetchPlugins, fetchProjects, fetchProjectServices, fetchServiceReleases } from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import type { CreateReleasePayload, PluginRecord, Project, Release, Service } from '@/api'
import { ApiError } from '@/api/request'

export type ReleaseRow = Release & {
  project_id: string
  project_name: string
  service_name: string
}

type ServiceOptionRecord = Service & {
  project_name: string
}

type ReleaseTimelineBucket = {
  date: string
  label: string
  count: number
  completed: number
  failed: number
  items: ReleaseRow[]
}

function getReleaseStatusTagColor(status?: string) {
  switch (status) {
    case 'failed':
      return { color: '#fee2e2', textColor: '#b91c1c' }
    case 'completed':
      return { color: '#dcfce7', textColor: '#15803d' }
    default:
      return { color: '#dbeafe', textColor: '#1d4ed8' }
  }
}

function getTimelineDate(row: ReleaseRow) {
  const raw = row.created_at || ''
  const parsed = Date.parse(raw)
  if (!Number.isNaN(parsed)) {
    return new Date(parsed)
  }

  return null
}

function formatTimelineLabel(date: Date) {
  return date.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
  })
}

export function useReleaseListService() {
  const message = useMessage()
  const router = useRouter()
  const authStore = useAuthStore()
  const loading = ref(false)
  const releaseSubmitting = ref(false)
  const releaseModalOpen = ref(false)
  const rows = ref<ReleaseRow[]>([])
  const services = ref<ServiceOptionRecord[]>([])
  const plugins = ref<PluginRecord[]>([])
  const filters = reactive({
    keyword: '',
  })
  const releaseForm = reactive<CreateReleasePayload & { service_id: string }>({
    service_id: '',
    plugin_id: '',
    tag: '',
    target: 'main',
    name: '',
    notes: '',
  })

  const canCreateRelease = computed(() =>
    authStore.canAccess({ permissions: [permission.releaseWrite] }),
  )

  const serviceOptions = computed(() =>
    services.value.map(service => ({
      label: `${service.name} · ${service.project_name}`,
      value: service.id,
    })),
  )

  const releaserOptions = computed(() =>
    plugins.value
      .filter(plugin => plugin.type === 'releaser')
      .map(plugin => ({ label: plugin.name, value: plugin.id })),
  )

  const filteredRows = computed(() => {
    const keyword = filters.keyword.trim().toLowerCase()

    return rows.value.filter((row) => {
      return !keyword || [
        row.tag,
        row.target,
        row.name,
        row.notes,
        row.project_name,
        row.service_name,
        row.status,
      ].some(value => value?.toLowerCase().includes(keyword))
    })
  })

  const timelineBuckets = computed<ReleaseTimelineBucket[]>(() => {
    const groups = new Map<string, ReleaseTimelineBucket>()

    for (const row of filteredRows.value) {
      const date = getTimelineDate(row)
      const key = date ? date.toISOString().slice(0, 10) : 'undated'
      const existing = groups.get(key)

      if (existing) {
        existing.count += 1
        existing.completed += row.status === 'completed' ? 1 : 0
        existing.failed += row.status === 'failed' ? 1 : 0
        existing.items.push(row)
        continue
      }

      groups.set(key, {
        date: key,
        label: date ? formatTimelineLabel(date) : 'Undated',
        count: 1,
        completed: row.status === 'completed' ? 1 : 0,
        failed: row.status === 'failed' ? 1 : 0,
        items: [row],
      })
    }

    return [...groups.values()]
      .sort((left, right) => right.date.localeCompare(left.date))
      .slice(0, 10)
  })

  const timelineMaxCount = computed(() =>
    Math.max(...timelineBuckets.value.map(item => item.count), 1),
  )

  const openService = (row: ReleaseRow) => {
    router.push({
      name: 'service-details',
      params: {
        serviceId: row.service_id,
      },
    })
  }

  const resetFilters = () => {
    filters.keyword = ''
  }

  const resetReleaseForm = () => {
    releaseForm.service_id = serviceOptions.value[0]?.value || ''
    releaseForm.plugin_id = releaserOptions.value[0]?.value || ''
    releaseForm.tag = ''
    releaseForm.target = 'main'
    releaseForm.name = ''
    releaseForm.notes = ''
  }

  const openReleaseModal = () => {
    resetReleaseForm()
    releaseModalOpen.value = true
  }

  const columns: DataTableColumns<ReleaseRow> = [
    { title: 'Tag', key: 'tag' },
    { title: 'Service', key: 'service_name' },
    { title: 'Project', key: 'project_name' },
    { title: 'Target', key: 'target' },
    {
      title: 'Status',
      key: 'status',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: getReleaseStatusTagColor(row.status),
          },
          { default: () => row.status || 'pending' },
        ),
    },
    {
      title: 'Release',
      key: 'html_url',
      render: row => row.html_url
        ? h(
            'a',
            {
              href: row.html_url,
              target: '_blank',
              rel: 'noreferrer',
              class: 'text-[var(--app-accent)] hover:underline',
            },
            'Open',
          )
        : 'Not set',
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
          { default: () => 'View service' },
        ),
    },
  ]

  const loadReleases = async() => {
    loading.value = true
    try {
      const [projects, pluginRows] = await Promise.all([
        fetchProjects(),
        fetchPlugins(),
      ])
      const projectGroups = await Promise.all(
        projects.map(async (project: Project) => {
          const projectServices = await fetchProjectServices(project.id)
          const serviceRows = projectServices.map(service => ({
            ...service,
            project_name: project.name,
          }))

          const releaseRows = await Promise.all(
            projectServices.map(async (service: Service) => {
              const releases = await fetchServiceReleases(service.id)

              return releases.map((release: Release) => ({
                ...release,
                project_id: project.id,
                project_name: project.name,
                service_name: service.name,
              }))
            }),
          )

          return {
            services: serviceRows,
            releases: releaseRows.flat(),
          }
        }),
      )

      services.value = projectGroups.flatMap(group => group.services)
      plugins.value = pluginRows
      rows.value = projectGroups.flatMap(group => group.releases)
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load releases.')
    } finally {
      loading.value = false
    }
  }

  const submitRelease = async() => {
    if (!releaseForm.service_id || !releaseForm.plugin_id || !releaseForm.tag.trim() || !releaseForm.target.trim()) {
      message.warning('Complete the release form before submitting.')
      return
    }

    releaseSubmitting.value = true
    try {
      await createRelease(releaseForm.service_id, {
        plugin_id: releaseForm.plugin_id,
        tag: releaseForm.tag.trim(),
        target: releaseForm.target.trim(),
        name: releaseForm.name?.trim() || undefined,
        notes: releaseForm.notes?.trim() || undefined,
      })
      message.success('Release created successfully.')
      releaseModalOpen.value = false
      await loadReleases()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to create release.')
    } finally {
      releaseSubmitting.value = false
    }
  }

  onMounted(loadReleases)

  return {
    canCreateRelease,
    columns,
    filteredRows,
    filters,
    loadReleases,
    loading,
    openReleaseModal,
    openService,
    releaseForm,
    releaseModalOpen,
    releaseSubmitting,
    releaserOptions,
    resetFilters,
    serviceOptions,
    submitRelease,
    timelineBuckets,
    timelineMaxCount,
  }
}
