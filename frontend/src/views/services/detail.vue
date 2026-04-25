<script setup lang="ts">
import {
  NButton,
  NCard,
  NCheckbox,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSelect,
  NStatistic,
  NTag,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { permission } from '@/services/access/rbac'
import PageHeader from '@/components/page-header.vue'
import {
  createDeployment,
  createScaffoldRequest,
  fetchDeploymentById,
  fetchPlugins,
  fetchProjectById,
  fetchProjects,
  fetchProjectServices,
  fetchServiceDeployments,
  fetchServiceReleases,
  suggestServiceScaffold,
} from '@/services/api'
import { ApiError } from '@/services/request'
import { useAuthStore } from '@/stores/modules/auth'
import { getEnvironmentTagColor } from '@/theme/environment'
import type {
  CreateDeploymentPayload,
  CreateScaffoldRequestPayload,
  Deployment,
  PluginRecord,
  Project,
  Release,
  ScaffoldSuggestion,
  Service,
} from '@/services/api'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const serviceId = computed(() => route.params.serviceId as string)
const projectId = computed(() => project.value?.id || service.value?.project_id || '')

const loading = ref(false)
const deploymentSubmitting = ref(false)
const deploymentLogLoading = ref(false)
const scaffoldSuggestionLoading = ref(false)
const scaffoldSubmitting = ref(false)
const deploymentModalOpen = ref(false)
const deploymentLogModalOpen = ref(false)
const scaffoldModalOpen = ref(false)
const project = ref<Project | null>(null)
const service = ref<Service | null>(null)
const releases = ref<Release[]>([])
const deployments = ref<Deployment[]>([])
const plugins = ref<PluginRecord[]>([])
const selectedReleaseTag = ref<string | null>(null)
const selectedDeployment = ref<Deployment | null>(null)

const deploymentForm = reactive<CreateDeploymentPayload>({
  plugin_id: '',
  environment: 'dev',
  version: '',
})

const scaffoldForm = reactive<CreateScaffoldRequestPayload>({
  plugin_id: '',
  environment: 'dev',
  variables: {
    service_name: '',
    module_path: '',
    port: 8080,
    database: 'postgres',
    enable_logging: true,
  },
})

const scaffoldSuggestion = ref<ScaffoldSuggestion | null>(null)

const successfulReleases = computed(() =>
  releases.value.filter(item => item.status === 'completed').length,
)

const failedReleases = computed(() =>
  releases.value.filter(item => item.status === 'failed').length,
)

const successfulDeployments = computed(() =>
  deployments.value.filter(item => item.status === 'completed').length,
)

const failedDeployments = computed(() =>
  deployments.value.filter(item => item.status === 'failed').length,
)

const deployerOptions = computed(() =>
  plugins.value
    .filter(plugin => plugin.type === 'deployer')
    .map(plugin => ({ label: plugin.name, value: plugin.id })),
)

const scaffolderOptions = computed(() =>
  plugins.value
    .filter(plugin => plugin.type === 'scaffolder')
    .map(plugin => ({ label: plugin.name, value: plugin.id })),
)

const environmentSelectOptions = computed(() => {
  const values = project.value?.environments?.length ? project.value.environments : ['dev', 'staging', 'prod']
  return values.map(value => ({ label: value, value }))
})

const canCreateDeployment = computed(() =>
  authStore.canAccess({ permissions: [permission.deploymentWrite] }),
)

const canCreateScaffoldRequest = computed(() =>
  authStore.canAccess({ permissions: [permission.scaffoldRequestWrite] }),
)

const selectedRelease = computed(() =>
  releases.value.find(item => item.tag === selectedReleaseTag.value) || null,
)

const visibleDeployments = computed(() =>
  selectedReleaseTag.value
    ? deployments.value.filter(item => item.version === selectedReleaseTag.value)
    : deployments.value,
)

function resetDeploymentForm() {
  deploymentForm.plugin_id = deployerOptions.value[0]?.value || ''
  deploymentForm.environment = project.value?.environments?.[0] || 'dev'
  deploymentForm.version = selectedRelease.value?.tag || ''
}

function resetScaffoldForm() {
  scaffoldForm.plugin_id = scaffolderOptions.value[0]?.value || ''
  scaffoldForm.environment = project.value?.environments?.[0] || 'dev'
  scaffoldForm.variables.service_name = normalizeServiceName(service.value?.name || project.value?.name || 'new-service')
  scaffoldForm.variables.module_path = inferModulePath(service.value?.repo_url || '')
  scaffoldForm.variables.port = suggestPort(scaffoldForm.variables.service_name)
  scaffoldForm.variables.database = 'postgres'
  scaffoldForm.variables.enable_logging = true
  scaffoldSuggestion.value = null
}

function openScaffoldModal() {
  resetScaffoldForm()
  generateScaffoldSuggestion()
  scaffoldModalOpen.value = true
}

function normalizeServiceName(value: string) {
  return value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    || 'new-service'
}

function inferModulePath(repoUrl: string) {
  const trimmed = repoUrl.trim().replace(/\.git$/, '')
  if (!trimmed) return ''
  return trimmed
    .replace(/^https?:\/\//, '')
    .replace(/^git@/, '')
    .replace(':', '/')
}

function suggestPort(name: string) {
  const hash = [...name].reduce((total, char) => total + char.charCodeAt(0), 0)
  return 8000 + (hash % 1000)
}

async function generateScaffoldSuggestion() {
  scaffoldSuggestionLoading.value = true

  try {
    scaffoldSuggestion.value = await suggestServiceScaffold(serviceId.value, {
      service_name: service.value?.name || scaffoldForm.variables.service_name,
      project_name: project.value?.name || '',
      project_description: project.value?.description || '',
      repo_url: service.value?.repo_url || scaffoldForm.variables.module_path,
      environment: scaffoldForm.environment,
      environments: project.value?.environments || [],
    })
  } catch (error) {
    message.warning(error instanceof ApiError
      ? `Using browser fallback: ${error.message}`
      : 'Using browser fallback for scaffold suggestions.')
  } finally {
    scaffoldSuggestionLoading.value = false
  }
}

function applyScaffoldSuggestion() {
  if (!scaffoldSuggestion.value) return

  scaffoldForm.environment = scaffoldSuggestion.value.environment
  scaffoldForm.variables = { ...scaffoldSuggestion.value.variables }
}

function selectRelease(row: Release) {
  selectedReleaseTag.value = row.tag
}

function clearReleaseSelection() {
  selectedReleaseTag.value = null
}

function formatDeploymentOutput(value?: string) {
  return value?.trim() || 'No runner output recorded yet.'
}

async function openDeploymentLogs(row: Deployment) {
  selectedDeployment.value = row
  deploymentLogModalOpen.value = true
  deploymentLogLoading.value = true

  try {
    selectedDeployment.value = await fetchDeploymentById(row.id)
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load deployment logs.')
  } finally {
    deploymentLogLoading.value = false
  }
}

const releaseColumns = computed(() => {
  const columns = [
    { title: 'Tag', key: 'tag' },
    { title: 'Target', key: 'target' },
    {
      title: 'Status',
      key: 'status',
      render: (row: Release) =>
        h(
          NTag,
          {
            bordered: false,
            color: row.status === 'failed'
              ? { color: '#fee2e2', textColor: '#b91c1c' }
              : { color: '#dbeafe', textColor: '#1d4ed8' },
          },
          { default: () => row.status || 'pending' },
        ),
    },
    {
      title: 'Release',
      key: 'html_url',
      render: (row: Release) => row.html_url
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
      render: (row: Release) =>
        h(
          NButton,
          {
            size: 'small',
            onClick: (event: MouseEvent) => {
              event.stopPropagation()
              selectRelease(row)
            },
          },
          { default: () => 'View deployments' },
        ),
    },
  ]

  if (canCreateDeployment.value) {
    columns.push({
      title: 'Deploy',
      key: 'deploy',
      render: (row: Release) =>
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            ghost: true,
            onClick: (event: MouseEvent) => {
              event.stopPropagation()
              selectRelease(row)
              resetDeploymentForm()
              deploymentForm.version = row.tag
              deploymentModalOpen.value = true
            },
          },
          { default: () => 'Deploy' },
        ),
    })
  }

  return columns
})

const deploymentColumns = [
  {
    title: 'Environment',
    key: 'environment',
    render: (row: Deployment) =>
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
    render: (row: Deployment) =>
      h(
        NTag,
        {
          bordered: false,
          color: row.status === 'failed'
            ? { color: '#fee2e2', textColor: '#b91c1c' }
            : { color: '#dbeafe', textColor: '#1d4ed8' },
        },
        { default: () => row.status || 'pending' },
      ),
  },
  { title: 'Commit SHA', key: 'commit_sha', render: (row: Deployment) => row.commit_sha || 'Not set' },
  {
    title: 'Runner output',
    key: 'runner_output',
    render: (row: Deployment) =>
      h(
        NButton,
        {
          size: 'small',
          onClick: (event: MouseEvent) => {
            event.stopPropagation()
            openDeploymentLogs(row)
          },
        },
        { default: () => row.runner_output || row.runner_error ? 'View output' : 'View logs' },
      ),
  },
]

async function loadServiceDetails() {
  loading.value = true

  try {
    const [serviceContext, pluginData, serviceReleases, serviceDeployments] = await Promise.all([
      findServiceContext(),
      fetchPlugins(),
      fetchServiceReleases(serviceId.value),
      fetchServiceDeployments(serviceId.value, { limit: 10, sortBy: 'date', sortOrder: 'desc' }),
    ])

    if (!serviceContext) {
      message.warning('Service not found.')
      router.push({ name: 'services' })
      return
    }

    project.value = serviceContext.project
    service.value = serviceContext.service
    plugins.value = pluginData
    releases.value = serviceReleases
    deployments.value = serviceDeployments
    resetScaffoldForm()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load service details.')
  } finally {
    loading.value = false
  }
}

async function findServiceContext(): Promise<{ project: Project, service: Service } | null> {
  const routeProjectId = route.params.projectId as string | undefined

  if (routeProjectId) {
    const [projectData, serviceRows] = await Promise.all([
      fetchProjectById(routeProjectId),
      fetchProjectServices(routeProjectId),
    ])
    const matchedService = serviceRows.find(item => item.id === serviceId.value) || null
    return matchedService ? { project: projectData, service: matchedService } : null
  }

  const projects = await fetchProjects()
  const serviceGroups = await Promise.all(
    projects.map(async (projectData: Project) => {
      const serviceRows = await fetchProjectServices(projectData.id)
      const matchedService = serviceRows.find(item => item.id === serviceId.value) || null
      return matchedService ? { project: projectData, service: matchedService } : null
    }),
  )

  return serviceGroups.find(Boolean) || null
}

async function submitScaffoldRequest() {
  if (!scaffoldForm.plugin_id || !scaffoldForm.environment || !scaffoldForm.variables.service_name.trim()) {
    message.warning('Complete the scaffold request form before submitting.')
    return
  }

  scaffoldSubmitting.value = true

  try {
    if (!projectId.value) {
      message.warning('Project context is not available for this scaffold request.')
      return
    }

    await createScaffoldRequest(projectId.value, {
      ...scaffoldForm,
      variables: {
        ...scaffoldForm.variables,
        service_name: scaffoldForm.variables.service_name.trim(),
        module_path: scaffoldForm.variables.module_path.trim(),
      },
    })
    message.success('Scaffold request created successfully.')
    scaffoldModalOpen.value = false
    resetScaffoldForm()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create scaffold request.')
  } finally {
    scaffoldSubmitting.value = false
  }
}

async function submitDeployment() {
  if (!deploymentForm.plugin_id || !deploymentForm.environment || !deploymentForm.version.trim()) {
    message.warning('Complete the deployment form before submitting.')
    return
  }

  deploymentSubmitting.value = true

  try {
    await createDeployment(serviceId.value, {
      plugin_id: deploymentForm.plugin_id,
      environment: deploymentForm.environment,
      version: deploymentForm.version.trim(),
    })
    message.success('Deployment created successfully.')
    deploymentModalOpen.value = false
    await loadServiceDetails()
    selectedReleaseTag.value = deploymentForm.version.trim()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create deployment.')
  } finally {
    deploymentSubmitting.value = false
  }
}

onMounted(loadServiceDetails)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Services"
      :title="service ? service.name : 'Service details'"
      description="Inspect the selected service, review its repository, and scan the most recent release and deployment activity tied to it."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="router.push({ name: 'services' })">
          Back to services
        </NButton>
        <NButton
          v-if="canCreateScaffoldRequest"
          secondary
          @click="openScaffoldModal"
        >
          New scaffold request
        </NButton>
        <!-- <NButton
          secondary
          :disabled="!selectedRelease"
          @click="openDeploymentModal"
        >
          Deploy selected release
        </NButton> -->
        <NButton
          v-if="service?.repo_url"
          tag="a"
          :href="service.repo_url"
          target="_blank"
          rel="noreferrer"
        >
          Open repository
        </NButton>
      </div>
    </PageHeader>

    <div class="grid gap-4 md:grid-cols-5">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Releases" :value="releases.length" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Successful releases" :value="successfulReleases" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Failed releases" :value="failedReleases" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Successful deployments" :value="successfulDeployments" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Failed deployments" :value="failedDeployments" />
      </NCard>
    </div>

    <div class="mt-6 grid gap-6 xl:grid-cols-[0.85fr_1.15fr]">
      <div class="grid gap-6">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Service posture">
          <div class="grid gap-4 text-sm leading-6 text-[var(--app-text-muted)]">
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Project
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ project?.name || 'Unknown' }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Repository
              </p>
              <a
                v-if="service?.repo_url"
                :href="service.repo_url"
                target="_blank"
                rel="noreferrer"
                class="mt-1 inline-flex text-base font-semibold text-[var(--app-accent)] hover:underline"
              >
                {{ service.repo_url }}
              </a>
              <p v-else class="mt-1 text-base font-semibold text-[var(--app-text)]">
                Not set
              </p>
            </div>
          </div>
        </NCard>
      </div>

      <div class="grid gap-6">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Recent releases">
          <NDataTable
            :columns="releaseColumns"
            :data="releases"
            :loading="loading"
            :pagination="{ pageSize: 6 }"
            :bordered="false"
            :row-props="(row: Release) => ({
              class: 'cursor-pointer',
              onClick: () => selectRelease(row),
            })"
          />
        </NCard>

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
          <template #header>
            <div class="flex flex-wrap items-center justify-between gap-3">
              <span class="text-lg font-600 text-[var(--app-text)]">Deployments</span>
              <div class="flex flex-wrap items-center gap-2">
                <NTag
                  v-if="selectedRelease"
                  :bordered="false"
                  :color="{ color: '#dbeafe', textColor: '#1d4ed8' }"
                >
                  {{ selectedRelease.tag }}
                </NTag>
                <NButton v-if="selectedRelease" size="small" @click="clearReleaseSelection">
                  Show all
                </NButton>
              </div>
            </div>
          </template>
          <NDataTable
            :columns="deploymentColumns"
            :data="visibleDeployments"
            :loading="loading"
            :pagination="{ pageSize: 6 }"
            :bordered="false"
            :row-props="(row: Deployment) => ({
              class: 'cursor-pointer',
              onClick: () => openDeploymentLogs(row),
            })"
          />
        </NCard>
      </div>
    </div>

    <NModal
      v-if="canCreateDeployment"
      v-model:show="deploymentModalOpen"
      preset="card"
      title="New deployment"
      class="max-w-2xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Deployer plugin">
            <NSelect
              v-model:value="deploymentForm.plugin_id"
              :options="deployerOptions"
              placeholder="Select deployer"
            />
          </NFormItem>

          <NFormItem label="Environment">
            <NSelect
              v-model:value="deploymentForm.environment"
              :options="environmentSelectOptions"
              placeholder="Select environment"
            />
          </NFormItem>

          <NFormItem label="Version" class="md:col-span-2">
            <NInput
              v-model:value="deploymentForm.version"
              placeholder="v1.0.0"
              :readonly="Boolean(selectedRelease)"
            />
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="deploymentModalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="deploymentSubmitting" @click="submitDeployment">
            Create deployment
          </NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-if="canCreateScaffoldRequest"
      v-model:show="scaffoldModalOpen"
      preset="card"
      title="New scaffold request"
      class="max-w-3xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="mb-4 rounded-2xl border border-[var(--app-border)] bg-slate-50 p-4">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <p class="text-sm font-700 text-[var(--app-text)]">AI-assisted scaffold suggestions</p>
                <NTag v-if="scaffoldSuggestion" size="small" :bordered="false" type="success">
                  {{ scaffoldSuggestion.source }}
                </NTag>
              </div>
              <p class="mt-1 text-xs text-[var(--app-text-muted)]">Generated by the local AI service from {{ service?.name || 'this service' }} and {{ project?.name || 'this project' }}</p>
            </div>
            <div class="flex gap-2">
              <NButton size="small" :loading="scaffoldSuggestionLoading" @click="generateScaffoldSuggestion">
                Refresh
              </NButton>
              <NButton
                size="small"
                type="primary"
                ghost
                :disabled="!scaffoldSuggestion || scaffoldSuggestionLoading"
                @click="applyScaffoldSuggestion"
              >
                Apply
              </NButton>
            </div>
          </div>

          <div v-if="scaffoldSuggestion" class="mt-4 grid gap-3 text-sm md:grid-cols-2">
            <div class="rounded-xl bg-white p-3">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Suggested service</p>
              <p class="mt-1 font-semibold text-[var(--app-text)]">{{ scaffoldSuggestion.variables.service_name }}</p>
            </div>
            <div class="rounded-xl bg-white p-3">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Suggested port</p>
              <p class="mt-1 font-semibold text-[var(--app-text)]">{{ scaffoldSuggestion.variables.port }}</p>
            </div>
            <div class="rounded-xl bg-white p-3">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Environment</p>
              <p class="mt-1 font-semibold text-[var(--app-text)]">{{ scaffoldSuggestion.environment }}</p>
            </div>
            <div class="rounded-xl bg-white p-3">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Module path</p>
              <p class="mt-1 break-all font-semibold text-[var(--app-text)]">{{ scaffoldSuggestion.variables.module_path || 'Manual input needed' }}</p>
            </div>
            <div class="rounded-xl bg-white p-3 md:col-span-2">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Reasoning</p>
              <ul class="mt-2 grid gap-1 text-[var(--app-text-muted)]">
                <li v-for="item in scaffoldSuggestion.rationale" :key="item">{{ item }}</li>
              </ul>
            </div>
          </div>
        </div>

        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Scaffolder plugin">
            <NSelect
              v-model:value="scaffoldForm.plugin_id"
              :options="scaffolderOptions"
              placeholder="Select scaffolder"
            />
          </NFormItem>

          <NFormItem label="Environment">
            <NSelect
              v-model:value="scaffoldForm.environment"
              :options="environmentSelectOptions"
              placeholder="Select environment"
            />
          </NFormItem>

          <NFormItem label="Service name">
            <NInput v-model:value="scaffoldForm.variables.service_name" placeholder="payments-api" />
          </NFormItem>

          <NFormItem label="Module path">
            <NInput v-model:value="scaffoldForm.variables.module_path" placeholder="github.com/acme/payments-api" />
          </NFormItem>

          <NFormItem label="Port">
            <NInputNumber v-model:value="scaffoldForm.variables.port" class="w-full" :min="1" :max="65535" />
          </NFormItem>

          <NFormItem label="Database">
            <NInput v-model:value="scaffoldForm.variables.database" placeholder="postgres" />
          </NFormItem>

          <NFormItem label="Enable logging">
            <NCheckbox v-model:checked="scaffoldForm.variables.enable_logging">
              Enabled
            </NCheckbox>
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="scaffoldModalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="scaffoldSubmitting" @click="submitScaffoldRequest">
            Create scaffold request
          </NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-model:show="deploymentLogModalOpen"
      preset="card"
      title="Deployment runner output"
      class="max-w-4xl"
      :bordered="false"
      segmented
    >
      <div class="grid gap-4">
        <div class="grid gap-3 text-sm md:grid-cols-3">
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
              :color="selectedDeployment?.status === 'failed'
                ? { color: '#fee2e2', textColor: '#b91c1c' }
                : { color: '#dbeafe', textColor: '#1d4ed8' }"
            >
              {{ selectedDeployment?.status || 'pending' }}
            </NTag>
          </div>
        </div>

        <NCard size="small" title="Runner output" :loading="deploymentLogLoading">
          <pre class="max-h-80 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950 p-4 text-xs leading-5 text-slate-100">{{ formatDeploymentOutput(selectedDeployment?.runner_output) }}</pre>
        </NCard>

        <NCard size="small" title="Runner errors">
          <pre class="max-h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-rose-950 p-4 text-xs leading-5 text-rose-50">{{ formatDeploymentOutput(selectedDeployment?.runner_error) }}</pre>
        </NCard>
      </div>
    </NModal>
  </div>
</template>
