package handler

import (
	"net/http"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	teamUsecase "devhub-backend/internal/usecase/team"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type updateTeamRequest struct {
	Name         *string `json:"name"`
	OwnerContact *string `json:"owner_contact"`
}

type updateTeamResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OwnerContact string `json:"owner_contact"`
}

// @Summary		Update Team
// @Description	Update an existing team
// @Tags			Team
// @Accept			json
// @Produce		json
// @Param			request	body		updateTeamRequest													true	"Team update input"
// @Success		200		{object}	httpresponse.SuccessResponse{data=updateTeamResponse,metadata=nil}	"Team updated"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}								"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}								"Internal server error"
// @Router			/teams/{team} [patch]
func (h *teamHandler) UpdateTeam(c *gin.Context) {
	var input updateTeamRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	updatedTeam, err := h.teamUsecase.UpdateTeam(c.Request.Context(), teamUsecase.UpdateTeamInput{
		ID:           c.Param("team"),
		Name:         input.Name,
		OwnerContact: input.OwnerContact,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusOK, h.newUpdateTeamResponse(updatedTeam))
}

func (h *teamHandler) newUpdateTeamResponse(team *entity.Team) updateTeamResponse {
	if team == nil {
		return updateTeamResponse{}
	}

	return updateTeamResponse{
		ID:           team.ID.String(),
		Name:         team.Name,
		OwnerContact: team.OwnerContact,
	}
}
