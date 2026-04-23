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

type createTeamRequest struct {
	Name         string `json:"name" binding:"required"`
	OwnerContact string `json:"owner_contact" binding:"required"`
}

type createTeamResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OwnerContact string `json:"owner_contact"`
}

// @Summary		Create Team
// @Description	Create a new team
// @Tags			Team
// @Accept			json
// @Produce		json
// @Param			request	body		createTeamRequest													true	"Team creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createTeamResponse,metadata=nil}	"Team created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}								"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}								"Internal server error"
// @Router			/teams [post]
func (h *teamHandler) CreateTeam(c *gin.Context) {
	var input createTeamRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdTeam, err := h.teamUsecase.CreateTeam(c.Request.Context(), teamUsecase.CreateTeamInput{
		Name:         input.Name,
		OwnerContact: input.OwnerContact,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateTeamResponse(createdTeam))
}

func (h *teamHandler) newCreateTeamResponse(team *entity.Team) createTeamResponse {
	if team == nil {
		return createTeamResponse{}
	}

	return createTeamResponse{
		ID:           team.ID.String(),
		Name:         team.Name,
		OwnerContact: team.OwnerContact,
	}
}
