<script setup lang="ts">
import { NButton, NCard, NDataTable, NInput, NSelect, NTag, useMessage } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import PageHeader from '@/components/page-header.vue'
import { fetchPlugins } from '@/services/api'
import { ApiError } from '@/services/request'
import { getPluginTypeTagColor, pluginTypeOptions } from '@/theme/plugin'
import type { PluginRecord } from '@/services/api'

const message = useMessage()
const loading = ref(false)
const rows = ref<PluginRecord[]>([])
const filters = reactive({
  keyword: '',
  type: null as string | null,
  runtime: null as string | null,
  scope: null as string | null,
})

const columns = [
  { title: 'Name', key: 'name' },
  {
    title: 'Type',
    key: 'type',
    render: (row: PluginRecord) =>
      h(
        NTag,
        { bordered: false, color: getPluginTypeTagColor(row.type) },
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

const runtimeOptions = computed(() =>
  [...new Set(rows.value.map(row => row.runtime))]
    .filter(Boolean)
    .map(value => ({ label: value, value })),
)

const scopeOptions = computed(() =>
  [...new Set(rows.value.map(row => row.scope))]
    .filter(Boolean)
    .map(value => ({ label: value, value })),
)

const filteredRows = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()

  return rows.value.filter((row) => {
    const matchesKeyword = !keyword || [
      row.name,
      row.entrypoint,
      row.version,
    ].some(value => value?.toLowerCase().includes(keyword))

    const matchesType = !filters.type || row.type === filters.type
    const matchesRuntime = !filters.runtime || row.runtime === filters.runtime
    const matchesScope = !filters.scope || row.scope === filters.scope

    return matchesKeyword && matchesType && matchesRuntime && matchesScope
  })
})

function resetFilters() {
  filters.keyword = ''
  filters.type = null
  filters.runtime = null
  filters.scope = null
}

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
      <NButton @click="load">
        Refresh
      </NButton>
    </PageHeader>

    <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1.2fr_0.8fr_0.8fr_0.8fr_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by name, version, or entrypoint"
          clearable
        />
        <NSelect
          v-model:value="filters.type"
          :options="pluginTypeOptions"
          placeholder="Type"
          clearable
        />
        <NSelect
          v-model:value="filters.runtime"
          :options="runtimeOptions"
          placeholder="Runtime"
          clearable
        />
        <NSelect
          v-model:value="filters.scope"
          :options="scopeOptions"
          placeholder="Scope"
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
      />
    </NCard>
  </div>
</template>
