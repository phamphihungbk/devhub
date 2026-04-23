package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	approvalUsecase "devhub-backend/internal/usecase/approval"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type FindAllApprovalRequestsQuery struct {
	Limit     *int64  `form:"limit"`
	Offset    *int64  `form:"offset"`
	SortBy    *string `form:"sortBy"`
	SortOrder *string `form:"sortOrder"`
	Status    *string `form:"status"`
}

func (h *approvalHandler) FindAllApprovalRequests(c *gin.Context) {
	var query FindAllApprovalRequestsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	var (
		limit     = misc.ToPointer(int64(100))
		offset    = misc.ToPointer(int64(0))
		sortBy    = misc.ToPointer("date")
		sortOrder = misc.ToPointer(entity.SortOrderAsc)
		status    *entity.ApprovalRequestStatus
	)
	if query.Limit != nil {
		limit = query.Limit
	}
	if query.Offset != nil {
		offset = query.Offset
	}
	if query.SortBy != nil {
		sortBy = query.SortBy
	}
	if query.SortOrder != nil {
		if querySortOrder, err := entity.ParseSortOrder(misc.GetValue(query.SortOrder)); err == nil {
			sortOrder = misc.ToPointer(querySortOrder)
		}
	}
	if query.Status != nil {
		parsedStatus, err := new(entity.ApprovalRequestStatus).Parse(misc.GetValue(query.Status))
		if err != nil {
			httpresponse.Error(c, errs.NewBadRequestError("invalid approval request status", nil))
			return
		}
		status = &parsedStatus
	}

	requests, err := h.approvalUsecase.FindAllApprovalRequests(c.Request.Context(), approvalUsecase.FindAllApprovalRequestsInput{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Status:    status,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithMetadata(c, h.newFindAllApprovalRequestsResponse(requests.GetData()), httpresponse.PaginationMetadata{Pagination: requests.GetPagination()})
}

func (h *approvalHandler) newFindAllApprovalRequestsResponse(requests entity.ApprovalRequests) []approvalRequestResponse {
	if len(requests) == 0 {
		return []approvalRequestResponse{}
	}

	response := make([]approvalRequestResponse, 0, len(requests))
	for _, request := range requests {
		response = append(response, h.newApprovalRequestResponse(&request))
	}

	return response
}
