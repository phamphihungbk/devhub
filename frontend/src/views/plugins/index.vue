<script setup lang="ts">
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NSelect,
  NSwitch,
  NTag,
} from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { usePluginService } from '@/services/plugin'

const {
  canManagePlugins,
  columns,
  detailModalOpen,
  filteredRows,
  filters,
  form,
  formModalOpen,
  formTitle,
  loadPlugins,
  loading,
  openCreatePlugin,
  resetFilters,
  runtimeOptions,
  runtimeSelectOptions,
  saving,
  scopeOptions,
  scopeSelectOptions,
  selectedPlugin,
  submitPlugin,
  typeOptions,
} = usePluginService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Automation"
      title="Plugins"
      description="The automation registry for deployers, releasers, and scaffolders. Audit runtime, scope, and entrypoints without the warmer marketing-style treatment."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="loadPlugins">
          Refresh
        </NButton>
        <NButton
          v-if="canManagePlugins"
          type="primary"
          @click="openCreatePlugin"
        >
          New plugin
        </NButton>
      </div>
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

    <NModal
      v-model:show="formModalOpen"
      preset="card"
      :title="formTitle"
      class="max-w-3xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Name" required>
            <NInput v-model:value="form.name" placeholder="go-http-api-scaffolder" />
          </NFormItem>

          <NFormItem label="Version" required>
            <NInput v-model:value="form.version" placeholder="1.0.0" />
          </NFormItem>

          <NFormItem label="Type" required>
            <NSelect
              v-model:value="form.type"
              :options="typeOptions"
              placeholder="Select type"
            />
          </NFormItem>

          <NFormItem label="Runtime" required>
            <NSelect
              v-model:value="form.runtime"
              :options="runtimeSelectOptions"
              placeholder="Select runtime"
            />
          </NFormItem>

          <NFormItem label="Scope" required>
            <NSelect
              v-model:value="form.scope"
              :options="scopeSelectOptions"
              placeholder="Select scope"
            />
          </NFormItem>

          <NFormItem label="Enabled">
            <NSwitch v-model:value="form.enabled" />
          </NFormItem>

          <NFormItem label="Entrypoint" required class="md:col-span-2">
            <NInput
              v-model:value="form.entrypoint"
              placeholder="/app/plugins/scaffolders/go_http_api/run.py"
            />
          </NFormItem>

          <NFormItem label="Description" required class="md:col-span-2">
            <NInput
              v-model:value="form.description"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="Describe when operators should use this plugin."
            />
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="formModalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="saving" @click="submitPlugin">
            Save plugin
          </NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-model:show="detailModalOpen"
      preset="card"
      :title="selectedPlugin?.name || 'Plugin details'"
      class="max-w-2xl"
      :bordered="false"
      segmented
    >
      <div v-if="selectedPlugin" class="grid gap-4 text-sm">
        <div class="flex flex-wrap gap-2">
          <NTag :bordered="false">
            {{ selectedPlugin.type }}
          </NTag>
          <NTag :bordered="false">
            {{ selectedPlugin.runtime }}
          </NTag>
          <NTag :bordered="false" :type="selectedPlugin.enabled === false ? 'default' : 'success'">
            {{ selectedPlugin.enabled === false ? 'Disabled' : 'Enabled' }}
          </NTag>
        </div>

        <div class="grid gap-3 md:grid-cols-2">
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Version</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedPlugin.version }}</p>
          </div>
          <div>
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Scope</p>
            <p class="mt-1 font-semibold text-[var(--app-text)]">{{ selectedPlugin.scope }}</p>
          </div>
          <div class="md:col-span-2">
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Entrypoint</p>
            <p class="mt-1 break-all font-semibold text-[var(--app-text)]">{{ selectedPlugin.entrypoint }}</p>
          </div>
          <div class="md:col-span-2">
            <p class="text-xs uppercase tracking-[0.22em] text-[var(--app-accent)]">Description</p>
            <p class="mt-1 leading-6 text-[var(--app-text-muted)]">{{ selectedPlugin.description || 'No description provided.' }}</p>
          </div>
        </div>
      </div>
    </NModal>
  </div>
</template>
