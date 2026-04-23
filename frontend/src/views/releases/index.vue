<script setup lang="ts">
import { NButton, NCard, NDataTable, NEmpty, NInput, NTag, useMessage } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import { fetchProjects, fetchProjectServices, fetchServiceReleases } from '@/services/api'
import { ApiError } from '@/services/request'
import type { Project, Release, Service } from '@/services/api'

type ReleaseRow = Release & {
  project_id: string
  project_name: string
  service_name: string
}

type ReleaseTimelineBucket = {
  date: string
  label: string
  count: number
  completed: number
  failed: number
  items: ReleaseRow[]
}

const message = useMessage()
const router = useRouter()
const loading = ref(false)
const rows = ref<ReleaseRow[]>([])
const filters = reactive({
  keyword: '',
})

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

function openService(row: ReleaseRow) {
  router.push({
    name: 'service-details',
    params: {
      projectId: row.project_id,
      serviceId: row.service_id,
    },
  })
}

const columns = [
  { title: 'Tag', key: 'tag' },
  { title: 'Service', key: 'service_name' },
  { title: 'Project', key: 'project_name' },
  { title: 'Target', key: 'target' },
  {
    title: 'Status',
    key: 'status',
    render: (row: ReleaseRow) =>
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
    render: (row: ReleaseRow) => row.html_url
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
    render: (row: ReleaseRow) =>
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

async function load() {
  loading.value = true
  try {
    const projects = await fetchProjects()
    const releaseGroups = await Promise.all(
      projects.map(async (project: Project) => {
        const services = await fetchProjectServices(project.id)

        return Promise.all(
          services.map(async (service: Service) => {
            const releases = await fetchServiceReleases(service.id)

            return releases.map((release: Release) => ({
              ...release,
              project_id: project.id,
              project_name: project.name,
              service_name: service.name,
            }))
          }),
        )
      }),
    )

    rows.value = releaseGroups.flatMap(group => group.flat())
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load releases.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Delivery"
      title="Releases"
      description="Track release activity across services, scan the most recent tags, and use the timeline chart to see how release volume moved over time."
    >
      <NButton @click="load">
        Refresh
      </NButton>
    </PageHeader>

    <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Release timeline">
      <template v-if="timelineBuckets.length > 0">
        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
          <div
            v-for="bucket in timelineBuckets"
            :key="bucket.date"
            class="rounded-3xl border border-[var(--app-border)] bg-white/86 p-4"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="text-xs font-700 uppercase tracking-0.22em text-brand-700">{{ bucket.label }}</p>
                <p class="mt-2 text-2xl font-700 text-ink-900">{{ bucket.count }}</p>
                <p class="text-sm text-[var(--app-text-muted)]">release events</p>
              </div>
              <div class="flex gap-2">
                <NTag :bordered="false" :color="{ color: '#dcfce7', textColor: '#15803d' }">
                  {{ bucket.completed }} ok
                </NTag>
                <NTag :bordered="false" :color="{ color: '#fee2e2', textColor: '#b91c1c' }">
                  {{ bucket.failed }} failed
                </NTag>
              </div>
            </div>

            <div class="mt-4">
              <div class="h-2 rounded-full bg-slate-100">
                <div
                  class="h-2 rounded-full bg-[linear-gradient(90deg,#2563eb_0%,#0f766e_100%)] transition-all"
                  :style="{ width: `${(bucket.count / timelineMaxCount) * 100}%` }"
                />
              </div>
            </div>

            <div class="mt-4 space-y-2">
              <div
                v-for="item in bucket.items.slice(0, 3)"
                :key="item.id"
                class="rounded-2xl bg-slate-50 px-3 py-2"
              >
                <p class="text-sm font-600 text-ink-900">{{ item.tag }} · {{ item.service_name }}</p>
                <p class="text-xs text-[var(--app-text-muted)]">{{ item.project_name }} · {{ new Date(item.created_at).toLocaleTimeString() }}</p>
              </div>
            </div>
          </div>
        </div>
      </template>
      <NEmpty v-else description="No release activity to chart yet." />
    </NCard>

    <NCard class="mt-6 rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1fr_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by tag, service, project, target, notes, or status"
          clearable
        />
        <NButton @click="filters.keyword = ''">
          Reset
        </NButton>
      </div>

      <NDataTable
        :columns="columns"
        :data="filteredRows"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
        :row-props="(row: ReleaseRow) => ({
          class: 'cursor-pointer',
          onClick: () => openService(row),
        })"
      />
    </NCard>
  </div>
</template>
