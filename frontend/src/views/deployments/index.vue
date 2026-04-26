<script setup lang="ts">
import { NButton, NCard, NDataTable, NEmpty, NForm, NFormItem, NInput, NModal, NSelect, NStatistic, NTag } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useDeploymentListService } from '@/services/deployment'

const {
  canCreateDeployment,
  columns,
  completedCount,
  deployerOptions,
  deploymentEnvironmentOptions,
  deploymentForm,
  deploymentModalOpen,
  deploymentSubmitting,
  environmentOptions,
  failedCount,
  filteredRows,
  filters,
  formatRunnerText,
  getDeploymentStatusTagColor,
  handleDeploymentServiceChange,
  loadDeployments,
  loading,
  logLoading,
  logModalOpen,
  openDeploymentModal,
  openLogs,
  resetFilters,
  rows,
  runningCount,
  selectedDeployment,
  serviceOptions,
  statusOptions,
  submitDeployment,
} = useDeploymentListService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Delivery"
      title="Deployments"
      description="Track deployment activity across services, inspect environment rollout state, and open runner output for completed or failed jobs."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="loadDeployments">
          Refresh
        </NButton>
        <NButton v-if="canCreateDeployment" type="primary" @click="openDeploymentModal">
          Deploy
        </NButton>
      </div>
    </PageHeader>

    <div class="grid gap-4 md:grid-cols-4">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Deployments" :value="rows.length" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Running" :value="runningCount" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Completed" :value="completedCount" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Failed" :value="failedCount" />
      </NCard>
    </div>

    <NCard class="mt-6 rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1fr_180px_180px_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by project, service, environment, version, status, commit, or reference"
          clearable
        />
        <NSelect
          v-model:value="filters.environment"
          :options="environmentOptions"
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
          onClick: () => openLogs(row),
        })"
      />
    </NCard>

    <NModal
      v-if="canCreateDeployment"
      v-model:show="deploymentModalOpen"
      preset="card"
      title="New deployment"
      class="max-w-2xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Service">
            <NSelect
              v-model:value="deploymentForm.service_id"
              :options="serviceOptions"
              placeholder="Select service"
              filterable
              @update:value="handleDeploymentServiceChange"
            />
          </NFormItem>

          <NFormItem label="Deployer plugin">
            <NSelect
              v-model:value="deploymentForm.plugin_id"
              :options="deployerOptions"
              placeholder="Select deployer"
            />
          </NFormItem>

          <NFormItem label="Environment">
            <NSelect
              v-model:value="deploymentForm.environment"
              :options="deploymentEnvironmentOptions"
              placeholder="Select environment"
            />
          </NFormItem>

          <NFormItem label="Version">
            <NInput v-model:value="deploymentForm.version" placeholder="v1.0.0" />
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="deploymentModalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="deploymentSubmitting" @click="submitDeployment">
            Create deployment
          </NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-model:show="logModalOpen"
      preset="card"
      title="Deployment runner output"
      class="max-w-4xl"
      :bordered="false"
      segmented
    >
      <div class="grid gap-4">
        <div class="grid gap-3 text-sm md:grid-cols-4">
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Service</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedDeployment?.service_name || 'Unknown' }}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Environment</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedDeployment?.environment || 'Unknown' }}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Version</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedDeployment?.version || 'Unknown' }}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Status</p>
            <NTag
              class="mt-1"
              :bordered="false"
              :color="getDeploymentStatusTagColor(selectedDeployment?.status)"
            >
              {{ selectedDeployment?.status || 'pending' }}
            </NTag>
          </div>
        </div>

        <NCard size="small" title="Runner output" :loading="logLoading">
          <pre class="max-h-80 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950 p-4 text-xs leading-5 text-slate-100">{{ formatRunnerText(selectedDeployment?.runner_output) }}</pre>
        </NCard>

        <NCard size="small" title="Runner errors">
          <pre class="max-h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-rose-950 p-4 text-xs leading-5 text-rose-50">{{ formatRunnerText(selectedDeployment?.runner_error) }}</pre>
        </NCard>
      </div>
    </NModal>

    <NEmpty v-if="!loading && rows.length === 0" class="mt-6" description="No deployments have been recorded yet." />
  </div>
</template>
