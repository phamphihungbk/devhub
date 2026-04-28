<script setup lang="ts">
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NPopconfirm,
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
  createServiceDependency,
  deleteServiceDependency,
  fetchDeploymentById,
  fetchPlugins,
  fetchProjectById,
  fetchProjects,
  fetchServiceDependencies,
  fetchProjectServices,
  fetchServiceDeployments,
  fetchServiceReleases,
} from '@/api'
import { ApiError } from '@/api/request'
import { useAuthStore } from '@/stores/modules/auth'
import { getEnvironmentTagColor } from '@/theme/environment'
import type {
  CreateDeploymentPayload,
  Deployment,
  PluginRecord,
  Project,
  Release,
  Service,
  ServiceDependency,
} from '@/api'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const serviceId = computed(() => route.params.serviceId as string)

const loading = ref(false)
const deploymentSubmitting = ref(false)
const deploymentLogLoading = ref(false)
const deploymentModalOpen = ref(false)
const deploymentLogModalOpen = ref(false)
const dependencyModalOpen = ref(false)
const project = ref<Project | null>(null)
const service = ref<Service | null>(null)
const projectServices = ref<Service[]>([])
const dependencies = ref<ServiceDependency[]>([])
const releases = ref<Release[]>([])
const deployments = ref<Deployment[]>([])
const plugins = ref<PluginRecord[]>([])
const selectedReleaseTag = ref<string | null>(null)
const selectedDeployment = ref<Deployment | null>(null)
const dependencySubmitting = ref(false)

const dependencyForm = reactive({
  depends_on_service_id: '',
  type: 'http',
  protocol: 'http',
  port: null as number | null,
  path: '',
})

const deploymentForm = reactive<CreateDeploymentPayload>({
  plugin_id: '',
  environment: 'dev',
  version: '',
})

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
    .filter(plugin => plugin.type === 'deployer' && plugin.enabled !== false)
    .map(plugin => ({ label: plugin.name, value: plugin.id })),
)

const environmentSelectOptions = computed(() => {
  const values = project.value?.environments?.length ? project.value.environments : ['dev', 'staging', 'prod']
  return values.map(value => ({ label: value, value }))
})

const canCreateDeployment = computed(() =>
  authStore.canAccess({ permissions: [permission.deploymentWrite] }),
)

const canWireService = computed(() =>
  authStore.canAccess({ permissions: [permission.projectWrite] }),
)

const selectedRelease = computed(() =>
  releases.value.find(item => item.tag === selectedReleaseTag.value) || null,
)

const visibleDeployments = computed(() =>
  selectedReleaseTag.value
    ? deployments.value.filter(item => item.version === selectedReleaseTag.value)
    : deployments.value,
)

const dependencyOptions = computed(() => {
  const wiredServiceIds = new Set(dependencies.value.map(item => item.depends_on_service_id))

  return projectServices.value
    .filter(item => item.id !== serviceId.value && !wiredServiceIds.has(item.id))
    .map(item => ({ label: item.name, value: item.id }))
})

const dependencyTypeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'gRPC', value: 'grpc' },
  { label: 'Queue', value: 'queue' },
  { label: 'Database', value: 'database' },
]

const dependencyProtocolOptions = computed(() => {
  if (dependencyForm.type === 'grpc') {
    return [{ label: 'gRPC', value: 'grpc' }]
  }
  if (dependencyForm.type === 'queue' || dependencyForm.type === 'database') {
    return [
      { label: 'TCP', value: 'tcp' },
      { label: 'UDP', value: 'udp' },
    ]
  }

  return [
    { label: 'HTTP', value: 'http' },
    { label: 'HTTPS', value: 'https' },
  ]
})

function resetDeploymentForm() {
  deploymentForm.plugin_id = deployerOptions.value[0]?.value || ''
  deploymentForm.environment = project.value?.environments?.[0] || 'dev'
  deploymentForm.version = selectedRelease.value?.tag || ''
}

function resetDependencyForm() {
  dependencyForm.depends_on_service_id = dependencyOptions.value[0]?.value || ''
  dependencyForm.type = 'http'
  dependencyForm.protocol = 'http'
  dependencyForm.port = null
  dependencyForm.path = ''
}

function openDependencyModal() {
  resetDependencyForm()
  dependencyModalOpen.value = true
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

const dependencyColumns = computed(() => {
  const columns = [
    {
      title: 'Depends on',
      key: 'depends_on_service',
      render: (row: ServiceDependency) => row.depends_on_service?.name || row.depends_on_service_id,
    },
    {
      title: 'Type',
      key: 'type',
      render: (row: ServiceDependency) =>
        h(
          NTag,
          { bordered: false, type: row.type === 'database' ? 'warning' : 'info' },
          { default: () => row.type },
        ),
    },
    {
      title: 'Endpoint',
      key: 'endpoint',
      render: (row: ServiceDependency) => {
        const protocol = row.protocol || 'default'
        const port = row.port ? `:${row.port}` : ''
        const path = row.path || ''
        return `${protocol}${port}${path}`
      },
    },
  ]

  if (canWireService.value) {
    columns.push({
      title: 'Actions',
      key: 'actions',
      render: (row: ServiceDependency) =>
        h(
          NPopconfirm,
          {
            onPositiveClick: () => removeDependency(row),
          },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  ghost: true,
                  onClick: (event: MouseEvent) => event.stopPropagation(),
                },
                { default: () => 'Remove' },
              ),
            default: () => 'Remove this service dependency?',
          },
        ),
    })
  }

  return columns
})

async function loadServiceDetails() {
  loading.value = true

  try {
    const [serviceContext, pluginData, serviceReleases, serviceDeployments, serviceDependencies] = await Promise.all([
      findServiceContext(),
      fetchPlugins(),
      fetchServiceReleases(serviceId.value),
      fetchServiceDeployments(serviceId.value, { limit: 10, sortBy: 'date', sortOrder: 'desc' }),
      fetchServiceDependencies(serviceId.value),
    ])

    if (!serviceContext) {
      message.warning('Service not found.')
      router.push({ name: 'services' })
      return
    }

    project.value = serviceContext.project
    service.value = serviceContext.service
    projectServices.value = serviceContext.services
    plugins.value = pluginData
    releases.value = serviceReleases
    deployments.value = serviceDeployments
    dependencies.value = serviceDependencies
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load service details.')
  } finally {
    loading.value = false
  }
}

async function findServiceContext(): Promise<{ project: Project, service: Service, services: Service[] } | null> {
  const routeProjectId = route.params.projectId as string | undefined

  if (routeProjectId) {
    const [projectData, serviceRows] = await Promise.all([
      fetchProjectById(routeProjectId),
      fetchProjectServices(routeProjectId),
    ])
    const matchedService = serviceRows.find(item => item.id === serviceId.value) || null
    return matchedService ? { project: projectData, service: matchedService, services: serviceRows } : null
  }

  const projects = await fetchProjects()
  const serviceGroups = await Promise.all(
    projects.map(async (projectData: Project) => {
      const serviceRows = await fetchProjectServices(projectData.id)
      const matchedService = serviceRows.find(item => item.id === serviceId.value) || null
      return matchedService ? { project: projectData, service: matchedService, services: serviceRows } : null
    }),
  )

  return serviceGroups.find(Boolean) || null
}

async function submitDependency() {
  if (!dependencyForm.depends_on_service_id || !dependencyForm.type) {
    message.warning('Choose a target service and dependency type before wiring.')
    return
  }

  dependencySubmitting.value = true

  try {
    await createServiceDependency(serviceId.value, {
      depends_on_service_id: dependencyForm.depends_on_service_id,
      type: dependencyForm.type,
      protocol: dependencyForm.protocol,
      port: dependencyForm.port,
      path: dependencyForm.path.trim(),
      config: {},
    })
    message.success('Service dependency wired successfully.')
    dependencyModalOpen.value = false
    dependencies.value = await fetchServiceDependencies(serviceId.value)
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to wire service dependency.')
  } finally {
    dependencySubmitting.value = false
  }
}

async function removeDependency(row: ServiceDependency) {
  try {
    await deleteServiceDependency(serviceId.value, row.id)
    dependencies.value = dependencies.value.filter(item => item.id !== row.id)
    message.success('Service dependency removed.')
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to remove service dependency.')
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
        <!-- <NButton
          v-if="service?.repo_url"
          tag="a"
          :href="service.repo_url"
          target="_blank"
          rel="noreferrer"
        >
          Open repository
        </NButton> -->
        <NButton
          v-if="canWireService"
          type="primary"
          ghost
          @click="openDependencyModal"
        >
          Wire service
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

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
          <template #header>
            <div class="flex flex-wrap items-center justify-between gap-3">
              <span class="text-lg font-600 text-[var(--app-text)]">Service wiring</span>
              <NButton
                v-if="canWireService"
                size="small"
                type="primary"
                ghost
                :disabled="dependencyOptions.length === 0"
                @click="openDependencyModal"
              >
                Add dependency
              </NButton>
            </div>
          </template>

          <NDataTable
            :columns="dependencyColumns"
            :data="dependencies"
            :loading="loading"
            :pagination="{ pageSize: 5 }"
            :bordered="false"
          />
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
      v-if="canWireService"
      v-model:show="dependencyModalOpen"
      preset="card"
      title="Wire service dependency"
      class="max-w-2xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Depends on">
            <NSelect
              v-model:value="dependencyForm.depends_on_service_id"
              :options="dependencyOptions"
              placeholder="Select service"
            />
          </NFormItem>

          <NFormItem label="Type">
            <NSelect
              v-model:value="dependencyForm.type"
              :options="dependencyTypeOptions"
              placeholder="Select type"
              @update:value="dependencyForm.protocol = dependencyProtocolOptions[0]?.value || ''"
            />
          </NFormItem>

          <NFormItem label="Protocol">
            <NSelect
              v-model:value="dependencyForm.protocol"
              :options="dependencyProtocolOptions"
              placeholder="Select protocol"
            />
          </NFormItem>

          <NFormItem label="Port">
            <NInputNumber
              v-model:value="dependencyForm.port"
              class="w-full"
              :min="1"
              :max="65535"
              placeholder="Optional"
            />
          </NFormItem>

          <NFormItem label="Path" class="md:col-span-2">
            <NInput
              v-model:value="dependencyForm.path"
              placeholder="/api/v1"
            />
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="dependencyModalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="dependencySubmitting" @click="submitDependency">
            Wire service
          </NButton>
        </div>
      </template>
    </NModal>

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
