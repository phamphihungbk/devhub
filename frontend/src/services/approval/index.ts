import { computed, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useRoute } from 'vue-router'
import { NButton, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { createApprovalDecision, fetchApprovalRequests } from '@/api'
import type { ApprovalRequestRecord } from '@/api'
import type {
  ApprovalAuditEventRecord,
  ApprovalRequestDetailRecord,
} from '@/api'
import { fetchApprovalRequestDetail } from '@/api'
import { ApiError } from '@/api/request'

const statusOptions = [
  { label: 'Pending', value: 'pending' },
  { label: 'Approved', value: 'approved' },
  { label: 'Rejected', value: 'rejected' },
  { label: 'Canceled', value: 'canceled' },
]

const getApprovalStatusTagColor = (status: string) => {
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

const formatScope = (row: ApprovalRequestRecord) => {
  return [row.project_id, row.service_id, row.environment].filter(Boolean).join(' / ') || 'Global'
}

const formatRequestedAt = (value: string) => {
  return new Date(value).toLocaleString()
}

export function useApprovalListService() {
  const message = useMessage()
  const router = useRouter()
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

  const loadApprovals = async() => {
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

  const openDecisionModal = (row: ApprovalRequestRecord, decision: 'approve' | 'reject') => {
    selectedRequest.value = row
    selectedDecision.value = decision
    decisionComment.value = ''
    decisionModalOpen.value = true
  }

  const openApprovalDetail = (row: ApprovalRequestRecord) => {
    router.push({
      name: 'approval-details',
      params: { approvalRequestId: row.id },
    })
  }

  const submitDecision = async() => {
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

  const columns: DataTableColumns<ApprovalRequestRecord> = [
    {
      title: 'Target',
      key: 'target',
      render: row => `${row.resource} / ${row.action}`,
    },
    {
      title: 'Scope',
      key: 'scope',
      render: row => formatScope(row),
    },
    {
      title: 'Status',
      key: 'status',
      render: row =>
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
      render: row => `${row.approved_count}/${row.required_approvals}`,
    },
    {
      title: 'Requested At',
      key: 'created_at',
      render: row => formatRequestedAt(row.created_at),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: row =>
        h(
          'div',
          { class: 'flex gap-2' },
          [
            h(
              NButton,
              {
                size: 'small',
                onClick: (event: MouseEvent) => {
                  event.stopPropagation()
                  openApprovalDetail(row)
                },
              },
              { default: () => 'View' },
            ),
            row.status === 'pending'
              ? h(
                  NButton,
                  {
                    size: 'small',
                    type: 'primary',
                    ghost: true,
                    loading: actingRequestId.value === row.id,
                    onClick: (event: MouseEvent) => {
                      event.stopPropagation()
                      openDecisionModal(row, 'approve')
                    },
                  },
                  { default: () => 'Approve' },
                )
              : null,
            row.status === 'pending'
              ? h(
                  NButton,
                  {
                    size: 'small',
                    type: 'error',
                    ghost: true,
                    loading: actingRequestId.value === row.id,
                    onClick: (event: MouseEvent) => {
                      event.stopPropagation()
                      openDecisionModal(row, 'reject')
                    },
                  },
                  { default: () => 'Reject' },
                )
              : null,
          ],
        ),
    },
  ]

  onMounted(loadApprovals)

  return {
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

function shortId(value: string) {
  return value ? `${value.slice(0, 8)}...${value.slice(-4)}` : 'Unknown'
}

export function useApprovalDetailService() {
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

  const loadApprovalDetail = async() => {
    loading.value = true
    try {
      detail.value = await fetchApprovalRequestDetail(approvalRequestId.value)
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load approval request.')
    } finally {
      loading.value = false
    }
  }

  const openApprovals = () => {
    router.push({ name: 'approvals' })
  }

  const openDecisionModal = (decision: 'approve' | 'reject') => {
    selectedDecision.value = decision
    decisionComment.value = ''
    decisionModalOpen.value = true
  }

  const submitDecision = async() => {
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
      await loadApprovalDetail()
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : `Unable to ${selectedDecision.value} approval request.`)
    } finally {
      acting.value = false
    }
  }

  onMounted(loadApprovalDetail)

  return {
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
  }
}
