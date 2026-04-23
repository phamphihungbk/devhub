package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	projectUsecase "devhub-backend/internal/usecase/project"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"time"

	"github.com/gin-gonic/gin"
)

type FindAllProjectsQuery struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02"`
	Limit     *int64     `form:"limit"`
	Offset    *int64     `form:"offset"`
	SortBy    *string    `form:"sortBy"`
	SortOrder *string    `form:"sortOrder"`
}

type findAllProjectsResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Project Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[development, production]"`
	Status       string   `json:"status" example:"active"`
	TeamID       string   `json:"team_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ScmProvider  string   `json:"scm_provider" example:"gitea"`
	CreatedBy    string   `json:"created_by" example:"Hung Pham"`
}

// @Summary		List Projects
// @Description	List all projects, filterable by date range and venue
// @Tags			Project
// @Produce		json
// @Param			startDate	query		string																									false	"Start date (format: 2006-01-02) (UTC+7)"
// @Param			endDate		query		string																									false	"End date (format: 2006-01-02) (UTC+7)"
// @Param			venue		query		string																									false	"Venue name (partial match)"
// @Param			limit		query		int64																									false	"Number of results to return (default: 100)"
// @Param			offset		query		int64																									false	"Number of results to skip (default: 0)"
// @Param			sortBy		query		string																									false	"Field to sort by (default: date) (options: date, name, venue)"
// @Param			sortOrder	query		string																									false	"Sort order (default: asc) (options: asc, desc)"
// @Success		200			{object}	httpresponse.SuccessResponse{data=[]findAllConcertsResponse,metadata=httpresponse.PaginationMetadata}	"List of concerts with pagination details"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}																	"Bad request"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}																	"Internal server error"
// @Router			/projects [get]
func (h *projectHandler) FindAllProjects(c *gin.Context) {
	var query FindAllProjectsQuery
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

	projects, err := h.projectUsecase.FindAllProjects(c.Request.Context(), projectUsecase.FindAllProjectsInput{
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

	httpresponse.SuccessWithMetadata(c, h.newFindAllProjectsResponse(projects.GetData()), httpresponse.PaginationMetadata{Pagination: projects.GetPagination()})
}

func (h *projectHandler) newFindAllProjectsResponse(projects entity.Projects) []findAllProjectsResponse {
	if len(projects) == 0 {
		return []findAllProjectsResponse{}
	}

	response := make([]findAllProjectsResponse, 0, len(projects))
	for _, project := range projects {
		envs := make([]string, 0, len(project.Environments))

		for _, env := range project.Environments {
			envs = append(envs, env.String())
		}

		response = append(response, findAllProjectsResponse{
			ID:           project.ID.String(),
			Name:         project.Name,
			Description:  project.Description,
			Environments: envs,
			Status:       project.Status.String(),
			TeamID:       project.TeamID.String(),
			ScmProvider:  project.ScmProvider,
			CreatedBy:    project.CreatedByName,
		})
	}
	return response
}
