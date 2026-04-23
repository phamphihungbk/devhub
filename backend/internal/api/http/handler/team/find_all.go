package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	teamUsecase "devhub-backend/internal/usecase/team"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type FindAllTeamsQuery struct {
	Limit     *int64  `form:"limit"`
	Offset    *int64  `form:"offset"`
	SortBy    *string `form:"sortBy"`
	SortOrder *string `form:"sortOrder"`
}

type findAllTeamsResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OwnerContact string `json:"owner_contact"`
}

// @Summary		List Teams
// @Description	List all teams
// @Tags			Team
// @Produce		json
// @Param			limit		query		int64	false	"Number of results to return (default: 100)"
// @Param			offset		query		int64	false	"Number of results to skip (default: 0)"
// @Param			sortBy		query		string	false	"Field to sort by (default: date) (options: date, name)"
// @Param			sortOrder	query		string	false	"Sort order (default: asc) (options: asc, desc)"
// @Success		200			{object}	httpresponse.SuccessResponse{data=[]findAllTeamsResponse,metadata=httpresponse.PaginationMetadata}	"List of teams"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}																"Bad request"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}																"Internal server error"
// @Router			/teams [get]
func (h *teamHandler) FindAllTeams(c *gin.Context) {
	var query FindAllTeamsQuery
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
		if parsedSortOrder, err := entity.ParseSortOrder(misc.GetValue(query.SortOrder)); err == nil {
			sortOrder = misc.ToPointer(parsedSortOrder)
		}
	}

	teams, err := h.teamUsecase.FindAllTeams(c.Request.Context(), teamUsecase.FindAllTeamsInput{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithMetadata(c, h.newFindAllTeamsResponse(teams.GetData()), httpresponse.PaginationMetadata{Pagination: teams.GetPagination()})
}

func (h *teamHandler) newFindAllTeamsResponse(teams entity.Teams) []findAllTeamsResponse {
	if len(teams) == 0 {
		return []findAllTeamsResponse{}
	}

	response := make([]findAllTeamsResponse, 0, len(teams))
	for _, team := range teams {
		response = append(response, findAllTeamsResponse{
			ID:           team.ID.String(),
			Name:         team.Name,
			OwnerContact: team.OwnerContact,
		})
	}

	return response
}
