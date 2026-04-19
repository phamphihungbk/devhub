<script setup lang="ts">
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSelect,
  NSpace,
  NStatistic,
  NTag,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import {
  createDeployment,
  createScaffoldRequest,
  fetchPlugins,
  fetchProjectById,
  fetchProjectServices,
  fetchProjectScaffoldRequests,
  fetchServiceDeployments,
} from '@/services/api'
import { ApiError } from '@/services/request'
import type {
  CreateDeploymentPayload,
  CreateScaffoldRequestPayload,
  Deployment,
  PluginRecord,
  Project,
  ScaffoldRequestRecord,
  Service,
} from '@/services/api'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const projectId = computed(() => route.params.projectId as string)

const pageLoading = ref(false)
const deploymentSubmitting = ref(false)
const scaffoldSubmitting = ref(false)

const project = ref<Project | null>(null)
const services = ref<Service[]>([])
const plugins = ref<PluginRecord[]>([])
const deployments = ref<Deployment[]>([])
const scaffoldRequests = ref<ScaffoldRequestRecord[]>([])
const selectedServiceId = ref('')

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
    port: 8080,
    database: 'postgres',
    enable_logging: true,
  },
})

const environmentOptions = [
  { label: 'Development', value: 'dev' },
  { label: 'Staging', value: 'staging' },
  { label: 'Production', value: 'prod' },
]

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

const pluginNameById = computed(() =>
  Object.fromEntries(plugins.value.map(plugin => [plugin.id, plugin.name])),
)

const serviceOptions = computed(() =>
  services.value.map(service => ({ label: service.name, value: service.id })),
)

const deploymentColumns = [
  {
    title: 'Status',
    key: 'status',
    render: (row: Deployment) =>
      h(
        NTag,
        {
          bordered: false,
          color: { color: '#dbeafe', textColor: '#1d4ed8' },
        },
        { default: () => row.status },
      ),
  },
  { title: 'Environment', key: 'environment' },
  { title: 'Version', key: 'version' },
  {
    title: 'Plugin',
    key: 'plugin_id',
    render: (row: Deployment) => pluginNameById.value[row.plugin_id] || row.plugin_id,
  },
  { title: 'Triggered By', key: 'triggered_by' },
]

const scaffoldColumns = [
  {
    title: 'Status',
    key: 'status',
    render: (row: ScaffoldRequestRecord) =>
      h(
        NTag,
        {
          bordered: false,
          color: { color: '#e2e8f0', textColor: '#334155' },
        },
        { default: () => row.status },
      ),
  },
  { title: 'Environment', key: 'environment' },
  {
    title: 'Plugin',
    key: 'plugin_id',
    render: (row: ScaffoldRequestRecord) => pluginNameById.value[row.plugin_id] || row.plugin_id,
  },
  {
    title: 'Service',
    key: 'service_name',
    render: (row: ScaffoldRequestRecord) => row.variables.service_name,
  },
  { title: 'Requested By', key: 'requested_by' },
]

async function loadProjectOperations() {
  pageLoading.value = true

  try {
    const [projectData, serviceData, pluginData, scaffoldData] = await Promise.all([
      fetchProjectById(projectId.value),
      fetchProjectServices(projectId.value),
      fetchPlugins(),
      fetchProjectScaffoldRequests(projectId.value),
    ])

    project.value = projectData
    services.value = serviceData
    plugins.value = pluginData
    scaffoldRequests.value = scaffoldData
    selectedServiceId.value = serviceData[0]?.id || ''
    deployments.value = selectedServiceId.value
      ? await fetchServiceDeployments(selectedServiceId.value)
      : []
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load project operations.')
  } finally {
    pageLoading.value = false
  }
}

function resetDeploymentForm() {
  deploymentForm.plugin_id = ''
  deploymentForm.environment = project.value?.environments?.[0] || 'dev'
  deploymentForm.version = ''
}

function resetScaffoldForm() {
  scaffoldForm.plugin_id = ''
  scaffoldForm.environment = project.value?.environments?.[0] || 'dev'
  scaffoldForm.variables.service_name = project.value?.name || ''
  scaffoldForm.variables.port = 8080
  scaffoldForm.variables.database = 'postgres'
  scaffoldForm.variables.enable_logging = true
}

async function submitDeployment() {
  if (!selectedServiceId.value || !deploymentForm.plugin_id || !deploymentForm.environment || !deploymentForm.version.trim()) {
    message.warning('Complete the deployment form before submitting.')
    return
  }

  deploymentSubmitting.value = true

  try {
    await createDeployment(selectedServiceId.value, {
      ...deploymentForm,
      version: deploymentForm.version.trim(),
    })
    message.success('Deployment created successfully.')
    resetDeploymentForm()
    deployments.value = await fetchServiceDeployments(selectedServiceId.value)
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create deployment.')
  } finally {
    deploymentSubmitting.value = false
  }
}

async function handleServiceChange(serviceId: string) {
  selectedServiceId.value = serviceId
  deployments.value = serviceId ? await fetchServiceDeployments(serviceId) : []
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
      },
    })
    message.success('Scaffold request created successfully.')
    resetScaffoldForm()
    scaffoldRequests.value = await fetchProjectScaffoldRequests(projectId.value)
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create scaffold request.')
  } finally {
    scaffoldSubmitting.value = false
  }
}

onMounted(async () => {
  await loadProjectOperations()
  resetDeploymentForm()
  resetScaffoldForm()
})
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Operations"
      :title="project ? `${project.name} operations` : 'Project operations'"
      description="Create project scaffold requests, deploy individual services, and track the most recent operational activity in one place."
    >
      <NSpace>
        <NButton @click="router.push({ name: 'projects' })">
          Back to projects
        </NButton>
        <NButton type="primary" secondary :loading="pageLoading" @click="loadProjectOperations">
          Refresh
        </NButton>
      </NSpace>
    </PageHeader>

    <div class="grid gap-4 md:grid-cols-4">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Deployments" :value="deployments.length" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Services" :value="services.length" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Scaffold requests" :value="scaffoldRequests.length" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Available environments" :value="project?.environments.length || 0" />
      </NCard>
    </div>

    <div class="mt-6 grid gap-6 xl:grid-cols-2">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="New deployment">
        <NForm label-placement="top">
          <div class="grid gap-4 md:grid-cols-2">
            <NFormItem label="Service">
              <NSelect
                v-model:value="selectedServiceId"
                :options="serviceOptions"
                placeholder="Select service"
                @update:value="handleServiceChange"
              />
            </NFormItem>

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
                :options="environmentOptions"
                placeholder="Select environment"
              />
            </NFormItem>

            <NFormItem label="Version" class="md:col-span-2">
              <NInput v-model:value="deploymentForm.version" placeholder="v1.0.0" />
            </NFormItem>
          </div>
        </NForm>

        <div class="mt-4 flex justify-end">
          <NButton type="primary" :loading="deploymentSubmitting" @click="submitDeployment">
            Create deployment
          </NButton>
        </div>
      </NCard>

      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="New scaffold request">
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

        <div class="mt-4 flex justify-end">
          <NButton type="primary" :loading="scaffoldSubmitting" @click="submitScaffoldRequest">
            Create scaffold request
          </NButton>
        </div>
      </NCard>
    </div>

    <div class="mt-6 grid gap-6 xl:grid-cols-2">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Deployment activity">
        <NDataTable
          :columns="deploymentColumns"
          :data="deployments"
          :loading="pageLoading"
          :pagination="{ pageSize: 5 }"
          :bordered="false"
        />
      </NCard>

      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Scaffold activity">
        <NDataTable
          :columns="scaffoldColumns"
          :data="scaffoldRequests"
          :loading="pageLoading"
          :pagination="{ pageSize: 5 }"
          :bordered="false"
        />
      </NCard>
    </div>
  </div>
</template>
