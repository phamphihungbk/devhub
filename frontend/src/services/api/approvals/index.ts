import { api } from '@/services/request'
import { apiBaseURL } from '../constants'
import type {
  ApprovalRequestListQuery,
  ApprovalRequestRecord,
  CreateApprovalDecisionPayload,
  CreateApprovalDecisionResponse,
} from './types'

const baseURL = `${apiBaseURL.approvalRequests}/`

export function fetchApprovalRequests(query?: ApprovalRequestListQuery) {
  return api.get<ApprovalRequestRecord[]>(baseURL, query)
}

export function createApprovalDecision(approvalRequestId: string, payload: CreateApprovalDecisionPayload) {
  return api.post<CreateApprovalDecisionResponse>(`${apiBaseURL.approvalRequests}/${approvalRequestId}/decisions`, payload)
}

export type {
  ApprovalRequestListQuery,
  ApprovalRequestRecord,
  CreateApprovalDecisionPayload,
  CreateApprovalDecisionResponse,
} from './types'
