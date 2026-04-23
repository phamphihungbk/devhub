<script setup lang="ts">
import { NButton, NCard, NDataTable, NInput, useMessage } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import { fetchProjects, fetchProjectServices } from '@/services/api'
import { ApiError } from '@/services/request'
import type { Project, Service } from '@/services/api'

type ServiceRow = Service & {
  project_name: string
}

const message = useMessage()
const router = useRouter()
const loading = ref(false)
const rows = ref<ServiceRow[]>([])
const filters = reactive({
  keyword: '',
})

const filteredRows = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()

  return rows.value.filter((row) => {
    return !keyword || [
      row.name,
      row.repo_url,
      row.project_name,
    ].some(value => value?.toLowerCase().includes(keyword))
  })
})

function openService(row: ServiceRow) {
  router.push({
    name: 'service-details',
    params: {
      projectId: row.project_id,
      serviceId: row.id,
    },
  })
}

const columns = [
  { title: 'Service', key: 'name' },
  { title: 'Project', key: 'project_name' },
  {
    title: 'Repository',
    key: 'repo_url',
    render: (row: ServiceRow) =>
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
    render: (row: ServiceRow) =>
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

async function load() {
  loading.value = true
  try {
    const projects = await fetchProjects()
    const serviceGroups = await Promise.all(
      projects.map(async (project: Project) => {
        const services = await fetchProjectServices(project.id)
        return services.map((service) => ({
          ...service,
          project_name: project.name,
        }))
      }),
    )

    rows.value = serviceGroups.flat()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load services.')
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
      title="Services"
      description="Browse the service inventory across projects, inspect repository links, and jump directly into a service’s release and deployment history."
    >
      <NButton @click="load">
        Refresh
      </NButton>
    </PageHeader>

    <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1fr_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by service, project, or repository"
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
        :row-props="(row: ServiceRow) => ({
          class: 'cursor-pointer',
          onClick: () => openService(row),
        })"
      />
    </NCard>
  </div>
</template>
