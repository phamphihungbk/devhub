<script setup lang="ts">
import { NButton, NCard, NDataTable, NEmpty, NInput, NModal, NSelect, NStatistic, NTag, useMessage } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import { fetchDeploymentById, fetchProjects, fetchProjectServices, fetchServiceDeployments } from '@/services/api'
import { ApiError } from '@/services/request'
import { getEnvironmentTagColor } from '@/theme/environment'
import type { Deployment, Project, Service } from '@/services/api'

type DeploymentRow = Deployment & {
  project_id: string
  project_name: string
  service_name: string
}

const message = useMessage()
const router = useRouter()
const loading = ref(false)
const logLoading = ref(false)
const rows = ref<DeploymentRow[]>([])
const selectedDeployment = ref<DeploymentRow | null>(null)
const logModalOpen = ref(false)
const filters = reactive({
  keyword: '',
  environment: null as string | null,
  status: null as string | null,
})

const environmentOptions = computed(() =>
  [...new Set(rows.value.map(row => row.environment).filter(Boolean))]
    .map(value => ({ label: value, value })),
)

const statusOptions = computed(() =>
  [...new Set(rows.value.map(row => row.status).filter(Boolean))]
    .map(value => ({ label: value, value })),
)

const filteredRows = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()

  return rows.value.filter((row) => {
    const matchesKeyword = !keyword || [
      row.project_name,
      row.service_name,
      row.environment,
      row.version,
      row.status,
      row.commit_sha,
      row.external_ref,
    ].some(value => value?.toLowerCase().includes(keyword))
    const matchesEnvironment = !filters.environment || row.environment === filters.environment
    const matchesStatus = !filters.status || row.status === filters.status

    return matchesKeyword && matchesEnvironment && matchesStatus
  })
})

const runningCount = computed(() =>
  rows.value.filter(row => row.status === 'running').length,
)

const completedCount = computed(() =>
  rows.value.filter(row => row.status === 'completed').length,
)

const failedCount = computed(() =>
  rows.value.filter(row => row.status === 'failed').length,
)

function getDeploymentStatusTagColor(status?: string) {
  switch (status) {
    case 'completed':
      return { color: '#dcfce7', textColor: '#15803d' }
    case 'failed':
      return { color: '#fee2e2', textColor: '#b91c1c' }
    case 'running':
      return { color: '#fef3c7', textColor: '#b45309' }
    default:
      return { color: '#dbeafe', textColor: '#1d4ed8' }
  }
}

function formatDate(value?: string) {
  return value ? new Date(value).toLocaleString() : 'Not set'
}

function formatRunnerText(value?: string) {
  return value?.trim() || 'No runner output recorded yet.'
}

function openService(row: DeploymentRow) {
  router.push({
    name: 'service-details',
    params: {
      serviceId: row.service_id,
    },
  })
}

async function openLogs(row: DeploymentRow) {
  selectedDeployment.value = row
  logModalOpen.value = true
  logLoading.value = true

  try {
    const deployment = await fetchDeploymentById(row.id)
    selectedDeployment.value = {
      ...row,
      ...deployment,
    }
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load deployment output.')
  } finally {
    logLoading.value = false
  }
}

function resetFilters() {
  filters.keyword = ''
  filters.environment = null
  filters.status = null
}

const columns = [
  { title: 'Service', key: 'service_name' },
  { title: 'Project', key: 'project_name' },
  {
    title: 'Environment',
    key: 'environment',
    render: (row: DeploymentRow) =>
      h(
        NTag,
        {
          bordered: false,
          color: getEnvironmentTagColor(row.environment),
        },
        { default: () => row.environment },
      ),
  },
  { title: 'Version', key: 'version' },
  {
    title: 'Status',
    key: 'status',
    render: (row: DeploymentRow) =>
      h(
        NTag,
        {
          bordered: false,
          color: getDeploymentStatusTagColor(row.status),
        },
        { default: () => row.status || 'pending' },
      ),
  },
  { title: 'Commit SHA', key: 'commit_sha', render: (row: DeploymentRow) => row.commit_sha || 'Not set' },
  { title: 'Finished', key: 'finished_at', render: (row: DeploymentRow) => formatDate(row.finished_at) },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: DeploymentRow) =>
      h(
        'div',
        { class: 'flex gap-2' },
        [
          h(
            NButton,
            {
              size: 'small',
              onClick: (event: MouseEvent) => {
                event.stopPropagation()
                openLogs(row)
              },
            },
            { default: () => 'Output' },
          ),
          h(
            NButton,
            {
              size: 'small',
              onClick: (event: MouseEvent) => {
                event.stopPropagation()
                openService(row)
              },
            },
            { default: () => 'Service' },
          ),
        ],
      ),
  },
]

async function load() {
  loading.value = true
  try {
    const projects = await fetchProjects()
    const deploymentGroups = await Promise.all(
      projects.map(async (project: Project) => {
        const services = await fetchProjectServices(project.id)

        return Promise.all(
          services.map(async (service: Service) => {
            const deployments = await fetchServiceDeployments(service.id, { limit: 50, sortBy: 'date', sortOrder: 'desc' })

            return deployments.map((deployment: Deployment) => ({
              ...deployment,
              project_id: project.id,
              project_name: project.name,
              service_name: service.name,
            }))
          }),
        )
      }),
    )

    rows.value = deploymentGroups.flatMap(group => group.flat())
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load deployments.')
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
      title="Deployments"
      description="Track deployment activity across services, inspect environment rollout state, and open runner output for completed or failed jobs."
    >
      <NButton @click="load">
        Refresh
      </NButton>
    </PageHeader>

    <div class="grid gap-4 md:grid-cols-4">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Deployments" :value="rows.length" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Running" :value="runningCount" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Completed" :value="completedCount" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Failed" :value="failedCount" />
      </NCard>
    </div>

    <NCard class="mt-6 rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1fr_180px_180px_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by project, service, environment, version, status, commit, or reference"
          clearable
        />
        <NSelect
          v-model:value="filters.environment"
          :options="environmentOptions"
          placeholder="Environment"
          clearable
        />
        <NSelect
          v-model:value="filters.status"
          :options="statusOptions"
          placeholder="Status"
          clearable
        />
        <NButton @click="resetFilters">
          Reset
        </NButton>
      </div>

      <NDataTable
        :columns="columns"
        :data="filteredRows"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
        :row-props="(row: DeploymentRow) => ({
          class: 'cursor-pointer',
          onClick: () => openLogs(row),
        })"
      />
    </NCard>

    <NModal
      v-model:show="logModalOpen"
      preset="card"
      title="Deployment runner output"
      class="max-w-4xl"
      :bordered="false"
      segmented
    >
      <div class="grid gap-4">
        <div class="grid gap-3 text-sm md:grid-cols-4">
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Service</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedDeployment?.service_name || 'Unknown' }}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Environment</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedDeployment?.environment || 'Unknown' }}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Version</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedDeployment?.version || 'Unknown' }}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Status</p>
            <NTag
              class="mt-1"
              :bordered="false"
              :color="getDeploymentStatusTagColor(selectedDeployment?.status)"
            >
              {{ selectedDeployment?.status || 'pending' }}
            </NTag>
          </div>
        </div>

        <NCard size="small" title="Runner output" :loading="logLoading">
          <pre class="max-h-80 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950 p-4 text-xs leading-5 text-slate-100">{{ formatRunnerText(selectedDeployment?.runner_output) }}</pre>
        </NCard>

        <NCard size="small" title="Runner errors">
          <pre class="max-h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-rose-950 p-4 text-xs leading-5 text-rose-50">{{ formatRunnerText(selectedDeployment?.runner_error) }}</pre>
        </NCard>
      </div>
    </NModal>

    <NEmpty v-if="!loading && rows.length === 0" class="mt-6" description="No deployments have been recorded yet." />
  </div>
</template>
