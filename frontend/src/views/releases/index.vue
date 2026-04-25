<script setup lang="ts">
import { NButton, NCard, NDataTable, NEmpty, NForm, NFormItem, NInput, NModal, NSelect, NTag, useMessage } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { permission } from '@/services/access/rbac'
import PageHeader from '@/components/page-header.vue'
import { createRelease, fetchPlugins, fetchProjects, fetchProjectServices, fetchServiceReleases } from '@/services/api'
import { ApiError } from '@/services/request'
import { useAuthStore } from '@/stores/modules/auth'
import type { CreateReleasePayload, PluginRecord, Project, Release, Service } from '@/services/api'

type ReleaseRow = Release & {
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
      serviceId: row.service_id,
    },
  })
}

function resetReleaseForm() {
  releaseForm.service_id = serviceOptions.value[0]?.value || ''
  releaseForm.plugin_id = releaserOptions.value[0]?.value || ''
  releaseForm.tag = ''
  releaseForm.target = 'main'
  releaseForm.name = ''
  releaseForm.notes = ''
}

function openReleaseModal() {
  resetReleaseForm()
  releaseModalOpen.value = true
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

async function submitRelease() {
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
    await load()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create release.')
  } finally {
    releaseSubmitting.value = false
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
      <div class="flex flex-wrap gap-3">
        <NButton @click="load">
          Refresh
        </NButton>
        <NButton v-if="canCreateRelease" type="primary" @click="openReleaseModal">
          New release
        </NButton>
      </div>
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

    <NModal
      v-if="canCreateRelease"
      v-model:show="releaseModalOpen"
      preset="card"
      title="New release"
      class="max-w-2xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Service">
            <NSelect
              v-model:value="releaseForm.service_id"
              :options="serviceOptions"
              placeholder="Select service"
              filterable
            />
          </NFormItem>

          <NFormItem label="Releaser plugin">
            <NSelect
              v-model:value="releaseForm.plugin_id"
              :options="releaserOptions"
              placeholder="Select releaser"
            />
          </NFormItem>

          <NFormItem label="Target">
            <NInput v-model:value="releaseForm.target" placeholder="main" />
          </NFormItem>

          <NFormItem label="Tag">
            <NInput v-model:value="releaseForm.tag" placeholder="v1.0.0" />
          </NFormItem>

          <NFormItem label="Name">
            <NInput v-model:value="releaseForm.name" placeholder="Payment v1.0.0" />
          </NFormItem>

          <NFormItem label="Notes" class="md:col-span-2">
            <NInput
              v-model:value="releaseForm.notes"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="Optional release notes or rollout context."
            />
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="releaseModalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="releaseSubmitting" @click="submitRelease">
            Create release
          </NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>
