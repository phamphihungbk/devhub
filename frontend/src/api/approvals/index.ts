import { apiBaseURL } from '../constants'
import { api } from '../request'
import type {
  ApprovalRequestListQuery,
  ApprovalRequestDetailRecord,
  ApprovalAuditEventRecord,
  ApprovalDecisionRecord,
  ApprovalRequestRecord,
  CreateApprovalDecisionPayload,
  CreateApprovalDecisionResponse,
} from './types'

const baseURL = `${apiBaseURL.approvalRequests}/`

export function fetchApprovalRequests(query?: ApprovalRequestListQuery) {
  return api.get<ApprovalRequestRecord[]>(baseURL, query)
}

export function fetchApprovalRequestDetail(approvalRequestId: string) {
  return api.get<ApprovalRequestDetailRecord>(`${apiBaseURL.approvalRequests}/${approvalRequestId}`)
}

export function createApprovalDecision(approvalRequestId: string, payload: CreateApprovalDecisionPayload) {
  return api.post<CreateApprovalDecisionResponse>(`${apiBaseURL.approvalRequests}/${approvalRequestId}/decisions`, payload)
}

export type {
  ApprovalAuditEventRecord,
  ApprovalDecisionRecord,
  ApprovalRequestListQuery,
  ApprovalRequestDetailRecord,
  ApprovalRequestRecord,
  CreateApprovalDecisionPayload,
  CreateApprovalDecisionResponse,
} from './types'
