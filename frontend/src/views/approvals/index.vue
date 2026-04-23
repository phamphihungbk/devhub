<script setup lang="ts">
import { NButton, NCard, NDataTable, NForm, NFormItem, NInput, NModal, NSelect, NTag, useMessage } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import PageHeader from '@/components/page-header.vue'
import { createApprovalDecision, fetchApprovalRequests } from '@/services/api'
import { ApiError } from '@/services/request'
import type { ApprovalRequestRecord } from '@/services/api'

const message = useMessage()
const loading = ref(false)
const actingRequestId = ref('')
const rows = ref<ApprovalRequestRecord[]>([])
const decisionModalOpen = ref(false)
const selectedRequest = ref<ApprovalRequestRecord | null>(null)
const selectedDecision = ref<'approve' | 'reject'>('approve')
const decisionComment = ref('')
const filters = reactive({
  status: 'pending',
})

const statusOptions = [
  { label: 'Pending', value: 'pending' },
  { label: 'Approved', value: 'approved' },
  { label: 'Rejected', value: 'rejected' },
  { label: 'Canceled', value: 'canceled' },
] as const

function getApprovalStatusTagColor(status: string) {
  switch (status) {
    case 'approved':
      return { color: '#dcfce7', textColor: '#15803d' }
    case 'rejected':
      return { color: '#fee2e2', textColor: '#b91c1c' }
    case 'canceled':
      return { color: '#e5e7eb', textColor: '#4b5563' }
    default:
      return { color: '#fef3c7', textColor: '#b45309' }
  }
}

function formatScope(row: ApprovalRequestRecord) {
  return [row.project_id, row.service_id, row.environment].filter(Boolean).join(' / ') || 'Global'
}

function formatRequestedAt(value: string) {
  return new Date(value).toLocaleString()
}

async function load() {
  loading.value = true
  try {
    rows.value = await fetchApprovalRequests({
      status: filters.status || undefined,
      sortBy: 'date',
      sortOrder: 'desc',
      limit: 100,
      offset: 0,
    })
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load approval requests.')
  } finally {
    loading.value = false
  }
}

function openDecisionModal(row: ApprovalRequestRecord, decision: 'approve' | 'reject') {
  selectedRequest.value = row
  selectedDecision.value = decision
  decisionComment.value = ''
  decisionModalOpen.value = true
}

async function submitDecision() {
  if (!selectedRequest.value) return

  const comment = decisionComment.value.trim()
  if (!comment) {
    message.warning('Comment is required before submitting a decision.')
    return
  }

  actingRequestId.value = selectedRequest.value.id
  try {
    const response = await createApprovalDecision(selectedRequest.value.id, {
      decision: selectedDecision.value,
      comment,
    })
    rows.value = rows.value.map(item =>
      item.id === selectedRequest.value?.id ? response.approval_request : item,
    )
    message.success(`Approval request ${selectedDecision.value}d successfully.`)
    decisionModalOpen.value = false
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : `Unable to ${selectedDecision.value} approval request.`)
  } finally {
    actingRequestId.value = ''
  }
}

const pendingCount = computed(() =>
  rows.value.filter(row => row.status === 'pending').length,
)

const columns = [
  {
    title: 'Target',
    key: 'target',
    render: (row: ApprovalRequestRecord) => `${row.resource} / ${row.action}`,
  },
  {
    title: 'Scope',
    key: 'scope',
    render: (row: ApprovalRequestRecord) => formatScope(row),
  },
  {
    title: 'Status',
    key: 'status',
    render: (row: ApprovalRequestRecord) =>
      h(
        NTag,
        {
          bordered: false,
          color: getApprovalStatusTagColor(row.status),
        },
        { default: () => row.status },
      ),
  },
  {
    title: 'Progress',
    key: 'progress',
    render: (row: ApprovalRequestRecord) => `${row.approved_count}/${row.required_approvals}`,
  },
  {
    title: 'Requested At',
    key: 'created_at',
    render: (row: ApprovalRequestRecord) => formatRequestedAt(row.created_at),
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: ApprovalRequestRecord) =>
      row.status !== 'pending'
        ? 'Resolved'
        : h(
            'div',
            { class: 'flex gap-2' },
            [
              h(
                NButton,
                {
                  size: 'small',
                  type: 'primary',
                  ghost: true,
                  loading: actingRequestId.value === row.id,
                  onClick: () => openDecisionModal(row, 'approve'),
                },
                { default: () => 'Approve' },
              ),
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  ghost: true,
                  loading: actingRequestId.value === row.id,
                  onClick: () => openDecisionModal(row, 'reject'),
                },
                { default: () => 'Reject' },
              ),
            ],
          ),
  },
]

onMounted(load)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Approvals"
      title="Approval requests"
      description="Review pending control-plane actions and resolve the requests that are currently waiting on manual approval."
    >
      <div class="flex flex-wrap gap-3">
        <NButton type="primary" @click="load">
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
          @update:value="load"
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
