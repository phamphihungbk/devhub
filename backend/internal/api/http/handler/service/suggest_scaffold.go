package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	serviceUsecase "devhub-backend/internal/usecase/service"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type suggestScaffoldRequest struct {
	ServiceName        string   `json:"service_name"`
	ProjectName        string   `json:"project_name"`
	ProjectDescription string   `json:"project_description"`
	RepoURL            string   `json:"repo_url"`
	Environment        string   `json:"environment"`
	Environments       []string `json:"environments"`
}

type suggestScaffoldResponse struct {
	Source      string                          `json:"source" example:"local-heuristic-v1"`
	Environment string                          `json:"environment" example:"dev"`
	Variables   entity.ScaffoldRequestVariables `json:"variables"`
	Rationale   []string                        `json:"rationale"`
}

// @Summary		Suggest Scaffold Defaults
// @Description	Generate local AI-assisted scaffold defaults from service and project context
// @Tags			Service
// @Accept			json
// @Produce		json
// @Param			service	path		string													true	"Service ID"
// @Param			request	body		suggestScaffoldRequest									true	"Scaffold suggestion context"
// @Success		200		{object}	httpresponse.SuccessResponse{data=suggestScaffoldResponse,metadata=nil}	"Scaffold suggestion"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}					"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}					"Internal server error"
// @Router			/services/{service}/scaffold-suggestions [post]
func (h *serviceHandler) SuggestScaffold(c *gin.Context) {
	var input suggestScaffoldRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	suggestion, err := h.serviceUsecase.SuggestScaffold(c.Request.Context(), serviceUsecase.SuggestScaffoldInput{
		ServiceID:          c.Param("service"),
		ServiceName:        input.ServiceName,
		ProjectName:        input.ProjectName,
		ProjectDescription: input.ProjectDescription,
		RepoURL:            input.RepoURL,
		Environment:        input.Environment,
		Environments:       input.Environments,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, suggestScaffoldResponse{
		Source:      suggestion.Source,
		Environment: suggestion.Environment,
		Variables:   suggestion.Variables,
		Rationale:   suggestion.Rationale,
	})
}
