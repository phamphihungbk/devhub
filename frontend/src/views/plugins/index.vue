<script setup lang="ts">
import { NButton, NCard, NDataTable, NInput, NSelect } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { usePluginService } from '@/services/plugin'

const {
  resetFilters,
  filteredRows,
  scopeOptions,
  runtimeOptions,
  typeOptions,
  filters,
  columns,
  loadPlugins,
  loading,
} = usePluginService()

</script>

<template>
  <div>
    <PageHeader
      eyebrow="Automation"
      title="Plugins"
      description="The automation registry for deployers, releasers, and scaffolders. Audit runtime, scope, and entrypoints without the warmer marketing-style treatment."
    >
      <NButton @click="loadPlugins">
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
          :options="typeOptions"
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
