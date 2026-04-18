<script setup lang="ts">
import { NButton, NCard, NDataTable, NTag, useMessage } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import PageHeader from '@/components/page-header.vue'
import { fetchPlugins } from '@/services/api'
import { ApiError } from '@/services/request'
import type { PluginRecord } from '@/services/api'

const message = useMessage()
const loading = ref(false)
const rows = ref<PluginRecord[]>([])

const columns = [
  { title: 'Name', key: 'name' },
  {
    title: 'Type',
    key: 'type',
    render: (row: PluginRecord) =>
      h(
        NTag,
        { bordered: false, color: { color: '#dbeafe', textColor: '#1d4ed8' } },
        { default: () => row.type },
      ),
  },
  {
    title: 'Runtime',
    key: 'runtime',
    render: (row: PluginRecord) =>
      h(
        NTag,
        { bordered: false, color: { color: '#e2e8f0', textColor: '#334155' } },
        { default: () => row.runtime },
      ),
  },
  { title: 'Version', key: 'version' },
  { title: 'Scope', key: 'scope' },
  { title: 'Entrypoint', key: 'entrypoint' },
]

async function load() {
  loading.value = true
  try {
    rows.value = await fetchPlugins()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load plugins.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Automation"
      title="Plugins"
      description="The automation registry for deployers, releasers, and scaffolders. Audit runtime, scope, and entrypoints without the warmer marketing-style treatment."
    >
      <NButton type="primary" secondary @click="load">
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
