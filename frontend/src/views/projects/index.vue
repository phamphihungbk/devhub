<script setup lang="ts">
import { NButton, NCard, NDataTable, NInput, NSelect } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useProjectListService } from '@/services/project'

const {
  canCreateProject,
  columns,
  environmentSelectOptions,
  filteredRows,
  filters,
  loadProjects,
  loading,
  openProject,
  openProjectCreate,
  ownerTeamOptions,
  resetFilters,
  statusOptions,
} = useProjectListService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Registry"
      title="Projects"
      description="The service catalog for the platform. Track ownership, deployment targets, and the spaces where scaffolding and lifecycle actions will land."
    >
      <div class="flex flex-wrap gap-3">
        <NButton v-if="canCreateProject" type="primary" @click="openProjectCreate">
          New project
        </NButton>
        <NButton @click="loadProjects">
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
          :options="environmentSelectOptions"
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
        :row-props="row => ({
          class: 'cursor-pointer',
          onClick: () => openProject(row),
        })"
      />
    </NCard>
  </div>
</template>
