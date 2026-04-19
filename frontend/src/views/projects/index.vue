<script setup lang="ts">
import { NButton, NCard, NDataTable, NInput, NSelect, NTag, useMessage } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import { fetchProjects } from '@/services/api'
import { ApiError } from '@/services/request'
import { environmentOptions, getEnvironmentTagColor } from '@/theme/environment'
import type { Project } from '@/services/api'

const message = useMessage()
const router = useRouter()
const loading = ref(false)
const rows = ref<Project[]>([])
const filters = reactive({
  keyword: '',
  status: null as string | null,
  environment: null as string | null,
  ownerTeam: null as string | null,
})

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Active', value: 'active' },
  { label: 'Archived', value: 'archived' },
  { label: 'Deprecated', value: 'deprecated' },
]

const ownerTeamOptions = computed(() =>
  [...new Set(rows.value.map(row => row.owner_team))]
    .filter(Boolean)
    .map(value => ({ label: value, value })),
)

const filteredRows = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()

  return rows.value.filter((row) => {
    const matchesKeyword = !keyword || [
      row.name,
      row.description,
      row.owner_team,
      row.status,
    ].some(value => value?.toLowerCase().includes(keyword))

    const matchesStatus = !filters.status || row.status === filters.status
    const matchesEnvironment = !filters.environment || row.environments.includes(filters.environment)
    const matchesOwnerTeam = !filters.ownerTeam || row.owner_team === filters.ownerTeam

    return matchesKeyword && matchesStatus && matchesEnvironment && matchesOwnerTeam
  })
})

function openProject(row: Project) {
  router.push({ name: 'project-details', params: { projectId: row.id } })
}

function resetFilters() {
  filters.keyword = ''
  filters.status = null
  filters.environment = null
  filters.ownerTeam = null
}

const columns = [
  { title: 'Name', key: 'name' },
  {
    title: 'Status',
    key: 'status',
    render: (row: Project) =>
      h(
        NTag,
        {
          bordered: false,
          color: { color: '#dbeafe', textColor: '#1d4ed8' },
        },
        { default: () => row.status || 'Unknown' },
      ),
  },
  {
    title: 'Environments',
    key: 'environments',
    render: (row: Project) =>
      h(
        'div',
        { class: 'flex flex-wrap gap-2' },
        row.environments.map((value) =>
          h(
            NTag,
            { bordered: false, color: getEnvironmentTagColor(value) },
            { default: () => value },
          ),
        ),
      ),
  },
  { title: 'Owner Team', key: 'owner_team', render: (row: Project) => row.owner_team || 'Not set' },
  { title: 'Description', key: 'description' },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: Project) =>
      h(
        NButton,
        {
          size: 'small',
          secondary: false,
          onClick: (event: MouseEvent) => {
            event.stopPropagation()
            openProject(row)
          },
        },
        { default: () => 'View details' },
      ),
  },
]

async function load() {
  loading.value = true
  try {
    rows.value = await fetchProjects()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load projects.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Registry"
      title="Projects"
      description="The service catalog for the platform. Track ownership, deployment targets, and the spaces where scaffolding and lifecycle actions will land."
    >
      <div class="flex flex-wrap gap-3">
        <NButton type="primary" @click="router.push({ name: 'project-create' })">
          New project
        </NButton>
        <NButton @click="load">
          Refresh
        </NButton>
      </div>
    </PageHeader>

    <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1.3fr_0.8fr_0.8fr_0.8fr_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by name, description, owner team, or status"
          clearable
        />
        <NSelect
          v-model:value="filters.status"
          :options="statusOptions"
          placeholder="Status"
          clearable
        />
        <NSelect
          v-model:value="filters.environment"
          :options="environmentOptions"
          placeholder="Environment"
          clearable
        />
        <NSelect
          v-model:value="filters.ownerTeam"
          :options="ownerTeamOptions"
          placeholder="Owner team"
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
        :row-props="(row: Project) => ({
          class: 'cursor-pointer',
          onClick: () => openProject(row),
        })"
      />
    </NCard>
  </div>
</template>
