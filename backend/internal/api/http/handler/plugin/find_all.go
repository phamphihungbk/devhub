package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"time"

	"github.com/gin-gonic/gin"
)

type FindAllPluginsQuery struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02" time_location:"Asia/Bangkok"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02" time_location:"Asia/Bangkok"`
	Limit     *int64     `form:"limit"`
	Offset    *int64     `form:"offset"`
	SortBy    *string    `form:"sortBy"`
	SortOrder *string    `form:"sortOrder"`
}

type findAllPluginsResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string `json:"name" example:"Plugin Name"`
	Type        string `json:"type" example:"scaffolder"`
	Version     string `json:"version" example:"1.0.0"`
	Description string `json:"description" example:"Plugin Description"`
}

// @Summary		List Plugins
// @Description	List all plugins, filterable by date range and venue
// @Tags			Plugin
// @Produce		json
// @Param			startDate	query		string																									false	"Start date (format: 2006-01-02) (UTC+7)"
// @Param			endDate		query		string																									false	"End date (format: 2006-01-02) (UTC+7)"
// @Param			venue		query		string																									false	"Venue name (partial match)"
// @Param			limit		query		int64																									false	"Number of results to return (default: 100)"
// @Param			offset		query		int64																									false	"Number of results to skip (default: 0)"
// @Param			sortBy		query		string																									false	"Field to sort by (default: date) (options: date, name, venue)"
// @Param			sortOrder	query		string																									false	"Sort order (default: asc) (options: asc, desc)"
// @Success		200			{object}	httpresponse.SuccessResponse{data=[]findAllPluginsResponse,metadata=httpresponse.PaginationMetadata}	"List of plugins with pagination details"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}																	"Bad request"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}																	"Internal server error"
// @Router			/plugins [get]
func (h *pluginHandler) FindAllPlugins(c *gin.Context) {
	var query FindAllPluginsQuery
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

	plugins, err := h.pluginUsecase.FindAllPlugins(c.Request.Context(), pluginUsecase.FindAllPluginsInput{
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

	httpresponse.SuccessWithMetadata(c, h.newFindAllPluginsResponse(plugins.GetData()), httpresponse.PaginationMetadata{Pagination: plugins.GetPagination()})
}

func (h *pluginHandler) newFindAllPluginsResponse(plugins entity.Plugins) []findAllPluginsResponse {
	if len(plugins) == 0 {
		return []findAllPluginsResponse{}
	}

	response := make([]findAllPluginsResponse, 0, len(plugins))
	for _, plugin := range plugins {
		response = append(response, findAllPluginsResponse{
			ID:   plugin.ID.String(),
			Name: plugin.Name,
		})
	}
	return response
}
