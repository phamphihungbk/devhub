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
  NSelect,
  NStatistic,
  NTag,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import {
  permission,
} from '@/access/rbac'
import {
  createScaffoldRequest,
  fetchPlugins,
  fetchProjectById,
  fetchProjectScaffoldRequests,
  fetchProjectServices,
  fetchServiceDeployments,
  fetchServiceReleases,
  fetchTeams,
} from '@/services/api'
import { ApiError } from '@/services/request'
import { useAuthStore } from '@/stores/modules/auth'
import { environmentOptions, getEnvironmentTagColor } from '@/theme/environment'
import type {
  CreateScaffoldRequestPayload,
  Deployment,
  PluginRecord,
  Project,
  Release,
  ScaffoldRequestRecord,
  Service,
  TeamRecord,
} from '@/services/api'

type ServiceReleaseRow = Release & { service_name: string }
type ServiceDeploymentRow = Deployment & { service_name: string }

const route = useRoute()
const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const projectId = computed(() => route.params.projectId as string)

const loading = ref(false)
const scaffoldSubmitting = ref(false)
const scaffoldModalOpen = ref(false)

const project = ref<Project | null>(null)
const teams = ref<TeamRecord[]>([])
const services = ref<Service[]>([])
const releases = ref<ServiceReleaseRow[]>([])
const deployments = ref<ServiceDeploymentRow[]>([])
const scaffoldRequests = ref<ScaffoldRequestRecord[]>([])
const plugins = ref<PluginRecord[]>([])

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

const scaffolderOptions = computed(() =>
  plugins.value
    .filter(plugin => plugin.type === 'scaffolder')
    .map(plugin => ({ label: plugin.name, value: plugin.id })),
)

const canCreateScaffoldRequest = computed(() =>
  authStore.canAccess({ permissions: [permission.scaffoldRequestWrite] }),
)

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

const teamNameById = computed(() =>
  new Map(teams.value.map(team => [team.id, team.name])),
)

const teamOwnerContactById = computed(() =>
  new Map(teams.value.map(team => [team.id, team.owner_contact])),
)

const ownerTeamName = computed(() =>
  project.value?.team_id ? (teamNameById.value.get(project.value.team_id) || project.value.team_id) : 'Not set',
)

const ownerContact = computed(() =>
  project.value?.team_id ? (teamOwnerContactById.value.get(project.value.team_id) || 'Not set') : 'Not set',
)

function openService(row: Service) {
  router.push({
    name: 'service-details',
    params: {
      projectId: projectId.value,
      serviceId: row.id,
    },
  })
}

const serviceColumns = [
  { title: 'Service', key: 'name' },
  {
    title: 'Repository',
    key: 'repo_url',
    render: (row: Service) =>
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
    render: (row: Service) =>
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

const releaseColumns = [
  { title: 'Service', key: 'service_name' },
  { title: 'Tag', key: 'tag' },
  { title: 'Target', key: 'target' },
  {
    title: 'Status',
    key: 'status',
    render: (row: ServiceReleaseRow) =>
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
]

const deploymentColumns = [
  { title: 'Service', key: 'service_name' },
  {
    title: 'Environment',
    key: 'environment',
    render: (row: ServiceDeploymentRow) =>
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
    render: (row: ServiceDeploymentRow) =>
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
]

function resetScaffoldForm() {
  scaffoldForm.plugin_id = ''
  scaffoldForm.environment = project.value?.environments?.[0] || 'dev'
  scaffoldForm.variables.service_name = project.value?.name || ''
  scaffoldForm.variables.module_path = ''
  scaffoldForm.variables.port = 8080
  scaffoldForm.variables.database = 'postgres'
  scaffoldForm.variables.enable_logging = true
}

async function loadProjectDetails() {
  loading.value = true

  try {
    const [projectData, serviceData, pluginData, scaffoldData] = await Promise.all([
      fetchProjectById(projectId.value),
      fetchProjectServices(projectId.value),
      fetchPlugins(),
      fetchProjectScaffoldRequests(projectId.value),
    ])
    teams.value = await fetchTeams()

    project.value = projectData
    services.value = serviceData
    plugins.value = pluginData
    scaffoldRequests.value = scaffoldData

    const serviceHistories = await Promise.all(
      serviceData.map(async (service) => {
        const [serviceReleases, serviceDeployments] = await Promise.all([
          fetchServiceReleases(service.id),
          fetchServiceDeployments(service.id, { limit: 5, sortBy: 'date', sortOrder: 'desc' }),
        ])

        return {
          releases: serviceReleases.map((item) => ({ ...item, service_name: service.name })),
          deployments: serviceDeployments.map((item) => ({ ...item, service_name: service.name })),
        }
      }),
    )

    releases.value = serviceHistories.flatMap(item => item.releases).slice(0, 8)
    deployments.value = serviceHistories.flatMap(item => item.deployments).slice(0, 8)
    resetScaffoldForm()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load project details.')
  } finally {
    loading.value = false
  }
}

async function submitScaffoldRequest() {
  if (!scaffoldForm.plugin_id || !scaffoldForm.environment || !scaffoldForm.variables.service_name.trim()) {
    message.warning('Complete the scaffold request form before submitting.')
    return
  }

  scaffoldSubmitting.value = true

  try {
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
    scaffoldRequests.value = await fetchProjectScaffoldRequests(projectId.value)
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create scaffold request.')
  } finally {
    scaffoldSubmitting.value = false
  }
}

onMounted(loadProjectDetails)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Registry"
      :title="project ? project.name : 'Project details'"
      description="Review the services under this project and the most recent lifecycle activity before deciding what to scaffold, release, or deploy next."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="router.push({ name: 'projects' })">
          Back to projects
        </NButton>
        <NButton v-if="canCreateScaffoldRequest" type="primary" @click="scaffoldModalOpen = true">
          New scaffold request
        </NButton>
      </div>
    </PageHeader>

    <div class="grid gap-4 md:grid-cols-5">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Services" :value="services.length" />
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

    <div class="mt-6 grid gap-6 xl:grid-cols-[0.95fr_1.05fr]">
      <div class="grid gap-6">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Project posture">
          <div class="grid gap-4 text-sm leading-6 text-[var(--app-text-muted)] md:grid-cols-2">
            <div class="md:col-span-2">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Environments
              </p>
              <div class="mt-3 flex flex-wrap gap-2">
                <NTag
                  v-for="environment in project?.environments || []"
                  :key="environment"
                  :bordered="false"
                  :color="getEnvironmentTagColor(environment)"
                >
                  {{ environment }}
                </NTag>
              </div>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Status
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ project?.status || 'Unknown' }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Owner team
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ ownerTeamName }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                SCM provider
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ project?.scm_provider || 'Not set' }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">
                Owner contact
              </p>
              <p class="mt-1 text-base font-semibold text-[var(--app-text)]">
                {{ ownerContact }}
              </p>
            </div>
          </div>
          <p class="mt-4 text-sm leading-6 text-[var(--app-text-muted)]">
            {{ project?.description || 'No description provided yet.' }}
          </p>
        </NCard>

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Services">
          <NDataTable
            :columns="serviceColumns"
            :data="services"
            :loading="loading"
            :pagination="{ pageSize: 6 }"
            :bordered="false"
            :row-props="(row: Service) => ({
              class: 'cursor-pointer',
              onClick: () => openService(row),
            })"
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
          />
        </NCard>

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Recent deployments">
          <NDataTable
            :columns="deploymentColumns"
            :data="deployments"
            :loading="loading"
            :pagination="{ pageSize: 6 }"
            :bordered="false"
          />
        </NCard>
      </div>
    </div>

    <NModal
      v-if="canCreateScaffoldRequest"
      v-model:show="scaffoldModalOpen"
      preset="card"
      title="New scaffold request"
      class="max-w-2xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
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
              :options="environmentOptions"
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
            <NSelect
              v-model:value="scaffoldForm.variables.enable_logging"
              :options="[
                { label: 'Enabled', value: true },
                { label: 'Disabled', value: false },
              ]"
            />
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
  </div>
</template>
