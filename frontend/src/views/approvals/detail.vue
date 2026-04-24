<script setup lang="ts">
import {
  NButton,
  NCard,
  NDescriptions,
  NDescriptionsItem,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NStatistic,
  NTag,
  NTimeline,
  NTimelineItem,
  useMessage,
} from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import PageHeader from '@/components/page-header.vue'
import { createApprovalDecision, fetchApprovalRequestDetail } from '@/services/api'
import { ApiError } from '@/services/request'
import type {
  ApprovalAuditEventRecord,
  ApprovalRequestDetailRecord,
  ApprovalRequestRecord,
} from '@/services/api'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const loading = ref(false)
const acting = ref(false)
const detail = ref<ApprovalRequestDetailRecord | null>(null)
const decisionModalOpen = ref(false)
const selectedDecision = ref<'approve' | 'reject'>('approve')
const decisionComment = ref('')

const approvalRequestId = computed(() => route.params.approvalRequestId as string)
const request = computed(() => detail.value?.approval_request || null)
const decisions = computed(() => detail.value?.decisions || [])
const auditEvents = computed(() => detail.value?.audit_events || [])
const remainingApprovals = computed(() => {
  if (!request.value) return 0
  return Math.max(request.value.required_approvals - request.value.approved_count, 0)
})

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

function getDecisionTagColor(decision: string) {
  return decision === 'reject'
    ? { color: '#fee2e2', textColor: '#b91c1c' }
    : { color: '#dcfce7', textColor: '#15803d' }
}

function getTimelineType(event: ApprovalAuditEventRecord) {
  if (event.type === 'reject' || event.type === 'rejected') return 'error'
  if (event.type === 'approve' || event.type === 'approved' || event.type === 'resolved') return 'success'
  return 'info'
}

function formatDate(value?: string | null) {
  return value ? new Date(value).toLocaleString() : 'Not set'
}

function formatScope(row: ApprovalRequestRecord) {
  return [row.project_id, row.service_id, row.environment].filter(Boolean).join(' / ') || 'Global'
}

function shortId(value: string) {
  return value ? `${value.slice(0, 8)}...${value.slice(-4)}` : 'Unknown'
}

async function load() {
  loading.value = true
  try {
    detail.value = await fetchApprovalRequestDetail(approvalRequestId.value)
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : 'Unable to load approval request.')
  } finally {
    loading.value = false
  }
}

function openDecisionModal(decision: 'approve' | 'reject') {
  selectedDecision.value = decision
  decisionComment.value = ''
  decisionModalOpen.value = true
}

async function submitDecision() {
  const comment = decisionComment.value.trim()
  if (!comment) {
    message.warning('Comment is required before submitting a decision.')
    return
  }

  acting.value = true
  try {
    await createApprovalDecision(approvalRequestId.value, {
      decision: selectedDecision.value,
      comment,
    })
    message.success(`Approval request ${selectedDecision.value}d successfully.`)
    decisionModalOpen.value = false
    await load()
  } catch (error) {
    message.error(error instanceof ApiError ? error.message : `Unable to ${selectedDecision.value} approval request.`)
  } finally {
    acting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Approvals"
      :title="request ? `${request.resource} / ${request.action}` : 'Approval request'"
      description="Review the approval target, decision trail, and audit timeline before resolving the request."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="router.push({ name: 'approvals' })">
          Back to approvals
        </NButton>
        <NButton @click="load">
          Refresh
        </NButton>
        <NButton
          v-if="request?.status === 'pending'"
          type="primary"
          :loading="acting"
          @click="openDecisionModal('approve')"
        >
          Approve
        </NButton>
        <NButton
          v-if="request?.status === 'pending'"
          type="error"
          ghost
          :loading="acting"
          @click="openDecisionModal('reject')"
        >
          Reject
        </NButton>
      </div>
    </PageHeader>

    <div v-if="request" class="grid gap-4 md:grid-cols-4">
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Approved" :value="request.approved_count" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Rejected" :value="request.rejected_count" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Required" :value="request.required_approvals" />
      </NCard>
      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
        <NStatistic label="Remaining" :value="remainingApprovals" />
      </NCard>
    </div>

    <div v-if="request" class="mt-6 grid gap-6 xl:grid-cols-[0.95fr_1.05fr]">
      <div class="grid gap-6">
        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Request detail">
          <NDescriptions :column="1" label-placement="left" bordered>
            <NDescriptionsItem label="Status">
              <NTag :bordered="false" :color="getApprovalStatusTagColor(request.status)">
                {{ request.status }}
              </NTag>
            </NDescriptionsItem>
            <NDescriptionsItem label="Scope">
              {{ formatScope(request) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="Resource ID">
              {{ request.resource_id }}
            </NDescriptionsItem>
            <NDescriptionsItem label="Requested by">
              {{ request.requested_by }}
            </NDescriptionsItem>
            <NDescriptionsItem label="Created">
              {{ formatDate(request.created_at) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="Resolved">
              {{ formatDate(request.resolved_at) }}
            </NDescriptionsItem>
          </NDescriptions>
        </NCard>

        <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Decision history">
          <div v-if="decisions.length" class="grid gap-3">
            <div
              v-for="decision in decisions"
              :key="decision.id"
              class="rounded-2xl border border-[var(--app-border)] p-4"
            >
              <div class="flex flex-wrap items-center justify-between gap-3">
                <NTag :bordered="false" :color="getDecisionTagColor(decision.decision)">
                  {{ decision.decision }}
                </NTag>
                <span class="text-xs text-[var(--app-text-muted)]">{{ formatDate(decision.created_at) }}</span>
              </div>
              <p class="mt-3 text-sm text-[var(--app-text)]">
                {{ decision.comment || 'No comment recorded.' }}
              </p>
              <p class="mt-2 text-xs text-[var(--app-text-muted)]">
                Decided by {{ decision.decided_by }}
              </p>
            </div>
          </div>
          <NEmpty v-else description="No decisions recorded yet." />
        </NCard>
      </div>

      <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]" title="Audit timeline">
        <NTimeline v-if="auditEvents.length">
          <NTimelineItem
            v-for="event in auditEvents"
            :key="`${event.type}-${event.created_at}-${event.actor_id}`"
            :type="getTimelineType(event)"
            :title="event.summary"
            :time="formatDate(event.created_at)"
          >
            <div class="text-sm leading-6 text-[var(--app-text-muted)]">
              <p>Actor {{ shortId(event.actor_id) }}</p>
              <p v-if="event.comment" class="mt-1 text-[var(--app-text)]">
                {{ event.comment }}
              </p>
            </div>
          </NTimelineItem>
        </NTimeline>
        <NEmpty v-else description="No audit events recorded yet." />
      </NCard>
    </div>

    <NCard
      v-else
      class="mt-6 rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]"
    >
      <NEmpty :description="loading ? 'Loading approval request...' : 'Approval request not found.'" />
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
            :loading="acting"
            @click="submitDecision"
          >
            {{ selectedDecision === 'approve' ? 'Approve' : 'Reject' }}
          </NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>
