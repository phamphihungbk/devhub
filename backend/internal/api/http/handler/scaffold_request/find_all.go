package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"time"

	"github.com/gin-gonic/gin"
)

type FindAllScaffoldRequestsQuery struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02" time_location:"Asia/Bangkok"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02" time_location:"Asia/Bangkok"`
	Limit     *int64     `form:"limit"`
	Offset    *int64     `form:"offset"`
	SortBy    *string    `form:"sortBy"`
	SortOrder *string    `form:"sortOrder"`
}

type findAllScaffoldRequestsResponse struct {
	ID          string                 `json:"id" example:"ad5b0c1f-762a-4ab3-a3e9-50a9057c49f3"`
	Template    string                 `json:"template" example:"go-service"`
	ProjectID   string                 `json:"project_id" example:"1a221b2c-abb7-44c0-8a96-8e92638b2422"`
	Environment string                 `json:"environment" example:"dev"`
	Variables   map[string]interface{} `json:"variables" example:"{\"service_name\":\"payment-service\",\"port\":8080,\"database\":\"postgres\",\"enable_logging\":true}"`
}

// @Summary		List Scaffold Requests
// @Description	List all scaffold requests, filterable by date range and venue
// @Tags			ScaffoldRequest
// @Produce		json
// @Param			startDate	query		string																									false	"Start date (format: 2006-01-02) (UTC+7)"
// @Param			endDate		query		string																									false	"End date (format: 2006-01-02) (UTC+7)"
// @Param			venue		query		string																									false	"Venue name (partial match)"
// @Param			limit		query		int64																									false	"Number of results to return (default: 100)"
// @Param			offset		query		int64																									false	"Number of results to skip (default: 0)"
// @Param			sortBy		query		string																									false	"Field to sort by (default: date) (options: date, name, venue)"
// @Param			sortOrder	query		string																									false	"Sort order (default: asc) (options: asc, desc)"
// @Success		200			{object}	httpresponse.SuccessResponse{data=[]findAllScaffoldRequestsResponse,metadata=httpresponse.PaginationMetadata}	"List of scaffold requests with pagination details"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}																	"Bad request"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}																	"Internal server error"
// @Router			/projects/:project/scaffold_requests [get]
func (h *scaffoldRequestHandler) FindAllScaffoldRequests(c *gin.Context) {
	projectID := c.Param("project")

	var query FindAllScaffoldRequestsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	var (
		limit     = misc.ToPointer(int64(100))          // Default limit to 100
		offset    = misc.ToPointer(int64(0))            // Default offset to 0
		sortBy    = misc.ToPointer("date")              // Default sort by date
		sortOrder = misc.ToPointer(entity.SortOrderAsc) // Default sort order
		err       error
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
		querySortOrder, err := entity.ParseSortOrder(misc.GetValue(query.SortOrder))
		if err == nil {
			sortOrder = misc.ToPointer(querySortOrder)
		}
	}

	scaffoldRequests, err := h.scaffoldRequestUsecase.FindAllScaffoldRequests(c.Request.Context(), scaffoldRequestUsecase.FindAllScaffoldRequestsInput{
		ProjectID: projectID,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithMetadata(c, h.newFindAllScaffoldRequestsResponse(scaffoldRequests.GetData()), httpresponse.PaginationMetadata{Pagination: scaffoldRequests.GetPagination()})
}

func (h *scaffoldRequestHandler) newFindAllScaffoldRequestsResponse(scaffoldRequests entity.ScaffoldRequests) []findAllScaffoldRequestsResponse {
	if len(scaffoldRequests) == 0 {
		return []findAllScaffoldRequestsResponse{}
	}

	response := make([]findAllScaffoldRequestsResponse, 0, len(scaffoldRequests))

	for _, scaffoldRequest := range scaffoldRequests {
		response = append(response, findAllScaffoldRequestsResponse{
			ID:          scaffoldRequest.ID.String(),
			Template:    scaffoldRequest.Template,
			ProjectID:   scaffoldRequest.ProjectID.String(),
			Environment: string(scaffoldRequest.Environment),
			Variables:   scaffoldRequest.Variables,
		})
	}
	return response
}
