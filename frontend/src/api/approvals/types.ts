export interface ApprovalRequestListQuery {
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  status?: string
}

export interface ApprovalRequestRecord {
  id: string
  resource: string
  resource_name?: string
  action: string
  resource_id: string
  requested_by: string
  requested_by_name?: string
  scope?: string
  project_id?: string
  service_id?: string
  environment?: string
  status: string
  required_approvals: number
  approved_count: number
  rejected_count: number
  resolved_at?: string
  created_at: string
  updated_at: string
}

export interface ApprovalDecisionRecord {
  id: string
  approval_request_id: string
  decided_by: string
  decided_by_name?: string
  decision: string
  comment: string
  created_at: string
}

export interface ApprovalAuditEventRecord {
  type: string
  actor_id: string
  actor_name?: string
  summary: string
  comment?: string
  created_at: string
}

export interface ApprovalRequestDetailRecord {
  approval_request: ApprovalRequestRecord
  decisions: ApprovalDecisionRecord[]
  audit_events: ApprovalAuditEventRecord[]
}

export interface CreateApprovalDecisionPayload {
  decision: 'approve' | 'reject'
  comment?: string
}

export interface CreateApprovalDecisionResponse {
  approval_request: ApprovalRequestRecord
  decision: ApprovalDecisionRecord
}
