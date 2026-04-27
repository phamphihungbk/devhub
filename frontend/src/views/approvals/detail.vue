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
} from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useApprovalDetailService } from '@/services/approval'

const {
  acting,
  auditEvents,
  decisionComment,
  decisionModalOpen,
  decisions,
  formatDate,
  formatScope,
  getApprovalStatusTagColor,
  getDecisionTagColor,
  getTimelineType,
  loadApprovalDetail,
  loading,
  openApprovals,
  openDecisionModal,
  remainingApprovals,
  request,
  selectedDecision,
  shortId,
  submitDecision,
} = useApprovalDetailService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Approvals"
      :title="request ? `${request.resource_name || request.resource} / ${request.action}` : 'Approval request'"
      description="Review the approval target, decision trail, and audit timeline before resolving the request."
    >
      <div class="flex flex-wrap gap-3">
        <NButton @click="openApprovals">
          Back to approvals
        </NButton>
        <NButton @click="loadApprovalDetail">
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
              {{ request.scope || formatScope(request) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="Resource">
              {{ request.resource_name || request.resource_id }}
            </NDescriptionsItem>
            <NDescriptionsItem label="Requested by">
              {{ request.requested_by_name || request.requested_by }}
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
                Decided by {{ decision.decided_by_name || decision.decided_by }}
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
              <p>{{ event.actor_name || `Actor ${shortId(event.actor_id)}` }}</p>
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
