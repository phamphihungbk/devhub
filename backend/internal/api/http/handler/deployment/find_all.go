package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	deploymentUsecase "devhub-backend/internal/usecase/deployment"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"time"

	"github.com/gin-gonic/gin"
)

type FindAllDeploymentsQuery struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02" time_location:"Asia/Bangkok"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02" time_location:"Asia/Bangkok"`
	Limit     *int64     `form:"limit"`
	Offset    *int64     `form:"offset"`
	SortBy    *string    `form:"sortBy"`
	SortOrder *string    `form:"sortOrder"`
}

type findAllDeploymentsResponse struct {
	ID   string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" example:"Project Name"`
}

// @Summary		List Deployments
// @Description	List all deployments, filterable by date range and venue
// @Tags			Deployment
// @Produce		json
// @Param			startDate	query		string																									false	"Start date (format: 2006-01-02) (UTC+7)"
// @Param			endDate		query		string																									false	"End date (format: 2006-01-02) (UTC+7)"
// @Param			venue		query		string																									false	"Venue name (partial match)"
// @Param			limit		query		int64																									false	"Number of results to return (default: 100)"
// @Param			offset		query		int64																									false	"Number of results to skip (default: 0)"
// @Param			sortBy		query		string																									false	"Field to sort by (default: date) (options: date, name, venue)"
// @Param			sortOrder	query		string																									false	"Sort order (default: asc) (options: asc, desc)"
// @Success		200			{object}	httpresponse.SuccessResponse{data=[]findAllDeploymentsResponse,metadata=httpresponse.PaginationMetadata}	"List of deployments with pagination details"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}																	"Bad request"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}																	"Internal server error"
// @Router			/projects/:project/deployment [get]
func (h *deploymentHandler) FindAllDeployments(c *gin.Context) {
	var query FindAllDeploymentsQuery
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

	deployments, err := h.deploymentUsecase.FindAllDeployments(c.Request.Context(), deploymentUsecase.FindAllDeploymentsInput{
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

	httpresponse.SuccessWithMetadata(c, h.newFindAllDeploymentsResponse(deployments.GetData()), httpresponse.PaginationMetadata{Pagination: deployments.GetPagination()})
}

func (h *deploymentHandler) newFindAllDeploymentsResponse(deployments entity.Deployments) []findAllDeploymentsResponse {
	if len(deployments) == 0 {
		return []findAllDeploymentsResponse{}
	}

	response := make([]findAllDeploymentsResponse, 0, len(deployments))
	for _, deployment := range deployments {
		response = append(response, findAllDeploymentsResponse{
			ID:   deployment.ID.String(),
			Name: deployment.Name,
		})
	}
	return response
}
