<script setup lang="ts">
import { NButton, NCard, NDataTable, NTag, useMessage } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import PageHeader from '@/components/page-header.vue'
import { fetchUsers } from '@/services/api'
import { ApiError } from '@/services/request'
import { getRoleTagColor } from '@/theme/role'
import type { UserRecord } from '@/services/api'

const message = useMessage()
const loading = ref(false)
const rows = ref<UserRecord[]>([])

const columns = [
  { title: 'Name', key: 'name' },
  { title: 'Email', key: 'email' },
  {
    title: 'Role',
    key: 'role',
    render: (row: UserRecord) =>
      h(
        NTag,
        { bordered: false, color: getRoleTagColor(row.role) },
        { default: () => row.role },
      ),
  },
]

async function load() {
  loading.value = true
  try {
    rows.value = await fetchUsers()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load users.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Access"
      title="Users"
      description="The operator directory for the control plane, ready for role governance, invite flows, and access reviews."
    >
      <NButton @click="load">
        Refresh
      </NButton>
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
