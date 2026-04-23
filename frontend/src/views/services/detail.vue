<script setup lang="ts">
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NSelect,
  NStatistic,
  NTag,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { permission } from '@/access/rbac'
import PageHeader from '@/components/page-header.vue'
import {
  createDeployment,
  createRelease,
  fetchPlugins,
  fetchProjectById,
  fetchProjectServices,
  fetchServiceDeployments,
  fetchServiceReleases,
} from '@/services/api'
import { ApiError } from '@/services/request'
import { useAuthStore } from '@/stores/modules/auth'
import { environmentOptions, getEnvironmentTagColor } from '@/theme/environment'
import type {
  CreateDeploymentPayload,
  CreateReleasePayload,
  Deployment,
  PluginRecord,
  Project,
  Release,
  Service,
} from '@/services/api'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const projectId = computed(() => route.params.projectId as string)
const serviceId = computed(() => route.params.serviceId as string)

const loading = ref(false)
const deploymentSubmitting = ref(false)
const releaseSubmitting = ref(false)
const deploymentModalOpen = ref(false)
const releaseModalOpen = ref(false)
const project = ref<Project | null>(null)
const service = ref<Service | null>(null)
const releases = ref<Release[]>([])
const deployments = ref<Deployment[]>([])
const plugins = ref<PluginRecord[]>([])
const selectedReleaseTag = ref<string | null>(null)

const releaseForm = reactive<CreateReleasePayload>({
  plugin_id: '',
  tag: '',
  target: 'main',
  name: '',
  notes: '',
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

const releaserOptions = computed(() =>
  plugins.value
    .filter(plugin => plugin.type === 'releaser')
    .map(plugin => ({ label: plugin.name, value: plugin.id })),
)

const deployerOptions = computed(() =>
  plugins.value
    .filter(plugin => plugin.type === 'deployer')
    .map(plugin => ({ label: plugin.name, value: plugin.id })),
)

const canCreateRelease = computed(() =>
  authStore.canAccess({ permissions: [permission.releaseWrite] }),
)

const canCreateDeployment = computed(() =>
  authStore.canAccess({ permissions: [permission.deploymentWrite] }),
)

const selectedRelease = computed(() =>
  releases.value.find(item => item.tag === selectedReleaseTag.value) || null,
)

const visibleDeployments = computed(() =>
  selectedReleaseTag.value
    ? deployments.value.filter(item => item.version === selectedReleaseTag.value)
    : deployments.value,
)

function resetReleaseForm() {
  releaseForm.plugin_id = releaserOptions.value[0]?.value || ''
  releaseForm.tag = ''
  releaseForm.target = 'main'
  releaseForm.name = ''
  releaseForm.notes = ''
}

function resetDeploymentForm() {
  deploymentForm.plugin_id = deployerOptions.value[0]?.value || ''
  deploymentForm.environment = project.value?.environments?.[0] || 'dev'
  deploymentForm.version = selectedRelease.value?.tag || ''
}

function openReleaseModal() {
  resetReleaseForm()
  releaseModalOpen.value = true
}

function openDeploymentModal() {
  if (!selectedRelease.value) {
    message.warning('Select a release first so deployment can use its version.')
    return
  }

  resetDeploymentForm()
  deploymentModalOpen.value = true
}

function selectRelease(row: Release) {
  selectedReleaseTag.value = row.tag
}

function clearReleaseSelection() {
  selectedReleaseTag.value = null
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
]

async function loadServiceDetails() {
  loading.value = true

  try {
    const [projectData, serviceRows, pluginData, serviceReleases, serviceDeployments] = await Promise.all([
      fetchProjectById(projectId.value),
      fetchProjectServices(projectId.value),
      fetchPlugins(),
      fetchServiceReleases(serviceId.value),
      fetchServiceDeployments(serviceId.value, { limit: 10, sortBy: 'date', sortOrder: 'desc' }),
    ])

    const matchedService = serviceRows.find(item => item.id === serviceId.value) || null

    if (!matchedService) {
      message.warning('Service not found for this project.')
      router.push({ name: 'project-details', params: { projectId: projectId.value } })
      return
    }

    project.value = projectData
    service.value = matchedService
    plugins.value = pluginData
    releases.value = serviceReleases
    deployments.value = serviceDeployments
    resetReleaseForm()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load service details.')
  } finally {
    loading.value = false
  }
}

async function submitRelease() {
  if (!releaseForm.plugin_id || !releaseForm.tag.trim() || !releaseForm.target.trim()) {
    message.warning('Complete the release form before submitting.')
    return
  }

  releaseSubmitting.value = true

  try {
    const createdTag = releaseForm.tag.trim()

    await createRelease(serviceId.value, {
      plugin_id: releaseForm.plugin_id,
      tag: createdTag,
      target: releaseForm.target.trim(),
      name: releaseForm.name?.trim() || undefined,
      notes: releaseForm.notes?.trim() || undefined,
    })
    message.success('Release created successfully.')
    releaseModalOpen.value = false
    await loadServiceDetails()
    selectedReleaseTag.value = createdTag
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to create release.')
  } finally {
    releaseSubmitting.value = false
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
        <NButton @click="router.push({ name: 'project-details', params: { projectId } })">
          Back to project
        </NButton>
        <NButton
          v-if="canCreateRelease"
          type="primary"
          @click="openReleaseModal"
        >
          New release
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
          />
        </NCard>
      </div>
    </div>

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
              :options="environmentOptions"
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
  </div>
</template>
