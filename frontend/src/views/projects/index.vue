<script setup lang="ts">
import { NButton, NCard, NDataTable, NTag, useMessage } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import { fetchProjects } from '@/services/api'
import { ApiError } from '@/services/request'
import type { Project } from '@/services/api'

const message = useMessage()
const router = useRouter()
const loading = ref(false)
const rows = ref<Project[]>([])

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
            { bordered: false, color: { color: '#e2e8f0', textColor: '#334155' } },
            { default: () => value },
          ),
        ),
      ),
  },
  { title: 'Owner Team', key: 'owner_team', render: (row: Project) => row.owner_team || 'Not set' },
  { title: 'SCM Provider', key: 'scm_provider', render: (row: Project) => row.scm_provider || 'Not set' },
  { title: 'Owner Contact', key: 'owner_contact', render: (row: Project) => row.owner_contact || 'Not set' },
  { title: 'Description', key: 'description' },
  { title: 'Created By', key: 'created_by' },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: Project) =>
      h(
        NButton,
        {
          size: 'small',
          secondary: true,
          onClick: () => router.push({ name: 'project-operations', params: { projectId: row.id } }),
        },
        { default: () => 'Manage ops' },
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
      <NDataTable
        :columns="columns"
        :data="rows"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
      />
    </NCard>
  </div>
</template>
