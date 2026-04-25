<script setup lang="ts">
import { NButton, NCard, NDataTable, NEmpty, NForm, NFormItem, NInput, NModal, NSelect, NTag } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useReleaseListService } from '@/services/release'

const {
  canCreateRelease,
  columns,
  filteredRows,
  filters,
  loadReleases,
  loading,
  openReleaseModal,
  openService,
  releaseForm,
  releaseModalOpen,
  releaseSubmitting,
  releaserOptions,
  resetFilters,
  serviceOptions,
  submitRelease,
  timelineBuckets,
  timelineMaxCount,
} = useReleaseListService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Delivery"
      title="Releases"
      description="Track release activity across services, scan the most recent tags, and use the timeline chart to see how release volume moved over time."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="loadReleases">
          Refresh
        </NButton>
        <NButton v-if="canCreateRelease" type="primary" @click="openReleaseModal">
          New release
        </NButton>
      </div>
    </PageHeader>

    <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Release timeline">
      <template v-if="timelineBuckets.length > 0">
        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
          <div
            v-for="bucket in timelineBuckets"
            :key="bucket.date"
            class="rounded-3xl border border-[var(--app-border)] bg-white/86 p-4"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="text-xs font-700 uppercase tracking-0.22em text-brand-700">{{ bucket.label }}</p>
                <p class="mt-2 text-2xl font-700 text-ink-900">{{ bucket.count }}</p>
                <p class="text-sm text-[var(--app-text-muted)]">release events</p>
              </div>
              <div class="flex gap-2">
                <NTag :bordered="false" :color="{ color: '#dcfce7', textColor: '#15803d' }">
                  {{ bucket.completed }} ok
                </NTag>
                <NTag :bordered="false" :color="{ color: '#fee2e2', textColor: '#b91c1c' }">
                  {{ bucket.failed }} failed
                </NTag>
              </div>
            </div>

            <div class="mt-4">
              <div class="h-2 rounded-full bg-slate-100">
                <div
                  class="h-2 rounded-full bg-[linear-gradient(90deg,#2563eb_0%,#0f766e_100%)] transition-all"
                  :style="{ width: `${(bucket.count / timelineMaxCount) * 100}%` }"
                />
              </div>
            </div>

            <div class="mt-4 space-y-2">
              <div
                v-for="item in bucket.items.slice(0, 3)"
                :key="item.id"
                class="rounded-2xl bg-slate-50 px-3 py-2"
              >
                <p class="text-sm font-600 text-ink-900">{{ item.tag }} · {{ item.service_name }}</p>
                <p class="text-xs text-[var(--app-text-muted)]">{{ item.project_name }} · {{ new Date(item.created_at).toLocaleTimeString() }}</p>
              </div>
            </div>
          </div>
        </div>
      </template>
      <NEmpty v-else description="No release activity to chart yet." />
    </NCard>

    <NCard class="mt-6 rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <div class="mb-5 grid gap-3 lg:grid-cols-[1fr_auto]">
        <NInput
          v-model:value="filters.keyword"
          placeholder="Filter by tag, service, project, target, notes, or status"
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

    <NModal
      v-if="canCreateRelease"
      v-model:show="releaseModalOpen"
      preset="card"
      title="New release"
      class="max-w-2xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="Service">
            <NSelect
              v-model:value="releaseForm.service_id"
              :options="serviceOptions"
              placeholder="Select service"
              filterable
            />
          </NFormItem>

          <NFormItem label="Releaser plugin">
            <NSelect
              v-model:value="releaseForm.plugin_id"
              :options="releaserOptions"
              placeholder="Select releaser"
            />
          </NFormItem>

          <NFormItem label="Target">
            <NInput v-model:value="releaseForm.target" placeholder="main" />
          </NFormItem>

          <NFormItem label="Tag">
            <NInput v-model:value="releaseForm.tag" placeholder="v1.0.0" />
          </NFormItem>

          <NFormItem label="Name">
            <NInput v-model:value="releaseForm.name" placeholder="Payment v1.0.0" />
          </NFormItem>

          <NFormItem label="Notes" class="md:col-span-2">
            <NInput
              v-model:value="releaseForm.notes"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="Optional release notes or rollout context."
            />
          </NFormItem>
        </div>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="releaseModalOpen = false">
            Cancel
          </NButton>
          <NButton type="primary" :loading="releaseSubmitting" @click="submitRelease">
            Create release
          </NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>
