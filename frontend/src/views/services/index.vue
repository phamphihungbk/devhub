<script setup lang="ts">
import { NButton, NCard, NDataTable, NInput } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useServiceListService } from '@/services/service'

const { columns, filteredRows, filters, loadServices, loading, openService, resetFilters } = useServiceListService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Registry"
      title="Services"
      description="Browse the service inventory across projects, inspect repository links, and jump directly into a service’s release and deployment history."
    >
      <NButton @click="loadServices">
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
        :row-props="row => ({
          class: 'cursor-pointer',
          onClick: () => openService(row),
        })"
      />
    </NCard>
  </div>
</template>
