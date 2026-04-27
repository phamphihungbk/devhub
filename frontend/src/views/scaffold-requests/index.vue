<script setup lang="ts">
import { NButton, NCard, NCheckbox, NDataTable, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect, NStatistic } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useScaffoldRequestListService } from '@/services/scaffold-request'

const {
  analyzePrompt,
  applySuggestion,
  canCreateScaffoldRequest,
  columns,
  completedCount,
  environmentOptions,
  failedCount,
  filterEnvironmentOptions,
  filteredRows,
  filters,
  form,
  handleProjectChange,
  loading,
  modalOpen,
  openModal,
  openProject,
  pendingCount,
  projectOptions,
  prompt,
  resetFilters,
  rows,
  runningCount,
  scaffolderOptions,
  statusOptions,
  submitting,
  submitScaffoldRequest,
  suggestion,
  suggestionLoading,
} = useScaffoldRequestListService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Scaffold requests"
      title="Scaffold Requests"
      description="Review scaffold activity across projects and create prompt-assisted requests before they enter approval and worker execution."
    >
      <NButton
        v-if="canCreateScaffoldRequest"
        type="primary"
        @click="openModal"
      >
        New scaffold request
      </NButton>
    </PageHeader>

    <div class="grid gap-4 md:grid-cols-4">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Pending" :value="pendingCount" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Running" :value="runningCount" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Completed" :value="completedCount" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Failed or rejected" :value="failedCount" />
      </NCard>
    </div>

    <NCard class="mt-6 rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1fr_180px_180px_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by service, project, module, database, status"
          clearable
        />
        <NSelect
          v-model:value="filters.environment"
          :options="filterEnvironmentOptions"
          placeholder="Environment"
          clearable
        />
        <NSelect
          v-model:value="filters.status"
          :options="statusOptions"
          placeholder="Status"
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

    <NModal
      v-if="canCreateScaffoldRequest"
      v-model:show="modalOpen"
      preset="card"
      title="New scaffold request"
      class="max-w-3xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Project">
            <NSelect
              v-model:value="form.project_id"
              :options="projectOptions"
              placeholder="Select project"
              filterable
              @update:value="handleProjectChange"
            />
          </NFormItem>

          <NFormItem label="Environment">
            <NSelect
              v-model:value="form.environment"
              :options="environmentOptions"
              placeholder="Select environment"
            />
          </NFormItem>

          <NFormItem label="Describe service" class="md:col-span-2">
            <NInput
              v-model:value="prompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="Create Go payment with mysql database and port 8000 with dev, prod and staging environment"
            />
          </NFormItem>

          <div class="md:col-span-2 flex justify-end gap-3">
            <NButton :loading="suggestionLoading" @click="analyzePrompt">
              Analyze prompt
            </NButton>
            <NButton
              type="primary"
              ghost
              :disabled="!suggestion"
              @click="applySuggestion"
            >
              Apply suggestion
            </NButton>
          </div>

          <div v-if="suggestion" class="md:col-span-2 grid gap-3 rounded-2xl border border-[var(--app-border)] bg-slate-50 p-4 md:grid-cols-2">
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Suggested plugin</p>
              <p class="mt-1 font-semibold text-[var(--app-text)]">{{ suggestion.plugin_name || 'Select manually' }}</p>
            </div>
            <div>
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Environment</p>
              <p class="mt-1 font-semibold text-[var(--app-text)]">{{ suggestion.environment }}</p>
            </div>
            <div class="md:col-span-2">
              <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Reasoning</p>
              <ul class="mt-2 grid gap-1 text-sm text-[var(--app-text-muted)]">
                <li v-for="item in suggestion.rationale" :key="item">{{ item }}</li>
              </ul>
            </div>
          </div>

          <NFormItem label="Scaffolder plugin">
            <NSelect
              v-model:value="form.plugin_id"
              :options="scaffolderOptions"
              placeholder="Select scaffolder"
              filterable
            />
          </NFormItem>

          <NFormItem label="Service name">
            <NInput v-model:value="form.variables.service_name" placeholder="payments-api" />
          </NFormItem>

          <NFormItem label="Module path">
            <NInput v-model:value="form.variables.module_path" placeholder="gitea.devhub.local/platform/payments-api" />
          </NFormItem>

          <NFormItem label="Port">
            <NInputNumber v-model:value="form.variables.port" class="w-full" :min="1" :max="65535" />
          </NFormItem>

          <NFormItem label="Database">
            <NInput v-model:value="form.variables.database" placeholder="postgres" />
          </NFormItem>

          <NFormItem label="Enable logging">
            <NCheckbox v-model:checked="form.variables.enable_logging">
              Enabled
            </NCheckbox>
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="modalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="submitting" @click="submitScaffoldRequest">
            Create scaffold request
          </NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>
