<script setup lang="ts">
import { NButton, NCard, NDataTable, NForm, NFormItem, NInput, NModal, NSelect } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useApprovalListService } from '@/services/approval'

const {
  actingRequestId,
  columns,
  decisionComment,
  decisionModalOpen,
  filters,
  loadApprovals,
  loading,
  openApprovalDetail,
  pendingCount,
  rows,
  selectedDecision,
  selectedRequest,
  statusOptions,
  submitDecision,
} = useApprovalListService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Approvals"
      title="Approval requests"
      description="Review pending control-plane actions and resolve the requests that are currently waiting on manual approval."
    >
      <div class="flex flex-wrap gap-3">
        <NButton type="primary" @click="loadApprovals">
          Refresh
        </NButton>
      </div>
    </PageHeader>

    <div class="mb-6 grid gap-4 md:grid-cols-3">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <p class="text-xs font-700 uppercase tracking-0.24em text-brand-700">Visible requests</p>
        <p class="mt-3 text-3xl font-700 text-ink-900">{{ rows.length }}</p>
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <p class="text-xs font-700 uppercase tracking-0.24em text-brand-700">Pending in view</p>
        <p class="mt-3 text-3xl font-700 text-ink-900">{{ pendingCount }}</p>
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <p class="text-xs font-700 uppercase tracking-0.24em text-brand-700">Status filter</p>
        <NSelect
          v-model:value="filters.status"
          class="mt-3"
          :options="statusOptions"
          @update:value="loadApprovals"
        />
      </NCard>
    </div>

    <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <NDataTable
        :columns="columns"
        :data="rows"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
        :row-props="row => ({
          class: 'cursor-pointer',
          onClick: () => openApprovalDetail(row),
        })"
      />
    </NCard>

    <NModal
      v-model:show="decisionModalOpen"
      preset="card"
      :title="selectedDecision === 'approve' ? 'Approve request' : 'Reject request'"
      class="max-w-xl"
      :bordered="false"
      segmented
    >
      <NForm label-placement="top">
        <NFormItem label="Comment" required>
          <NInput
            v-model:value="decisionComment"
            type="textarea"
            :autosize="{ minRows: 4, maxRows: 6 }"
            placeholder="Explain why you are approving or rejecting this request."
          />
        </NFormItem>
      </NForm>

      <template #action>
        <div class="flex justify-end gap-3">
          <NButton @click="decisionModalOpen = false">
            Cancel
          </NButton>
          <NButton
            :type="selectedDecision === 'approve' ? 'primary' : 'error'"
            :loading="Boolean(selectedRequest && actingRequestId === selectedRequest.id)"
            @click="submitDecision"
          >
            {{ selectedDecision === 'approve' ? 'Approve' : 'Reject' }}
          </NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>
