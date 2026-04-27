package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type suggestScaffoldRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type suggestScaffoldRequestResponse struct {
	Source       string                          `json:"source" example:"local-prompt-heuristic-v1"`
	PluginID     string                          `json:"plugin_id" example:"72bd5b8f-54b3-442a-b54f-685643f6d46e"`
	PluginName   string                          `json:"plugin_name" example:"Go HTTP API Scaffolder"`
	Confidence   float64                         `json:"confidence" example:"0.82"`
	Environment  string                          `json:"environment" example:"dev"`
	Environments []string                        `json:"environments" example:"dev,staging,prod"`
	Variables    entity.ScaffoldRequestVariables `json:"variables"`
	Rationale    []string                        `json:"rationale"`
}

// @Summary		Suggest Scaffold Request
// @Description	Analyze a user prompt and available scaffold plugins to suggest scaffold_request values
// @Tags			ScaffoldRequest
// @Accept			json
// @Produce		json
// @Param			project	path		string														true	"Project ID"
// @Param			request	body		suggestScaffoldRequest										true	"Scaffold suggestion prompt"
// @Success		200		{object}	httpresponse.SuccessResponse{data=suggestScaffoldRequestResponse,metadata=nil}	"Scaffold suggestion"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}						"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}						"Internal server error"
// @Router			/projects/{project}/scaffold-suggestions [post]
func (h *scaffoldRequestHandler) SuggestScaffoldRequest(c *gin.Context) {
	var input suggestScaffoldRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	suggestion, err := h.scaffoldRequestUsecase.SuggestScaffoldRequest(c.Request.Context(), scaffoldRequestUsecase.SuggestScaffoldRequestInput{
		ProjectID: c.Param("project"),
		Prompt:    input.Prompt,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, suggestScaffoldRequestResponse{
		Source:       suggestion.Source,
		PluginID:     suggestion.PluginID,
		PluginName:   suggestion.PluginName,
		Confidence:   suggestion.Confidence,
		Environment:  suggestion.Environment,
		Environments: suggestion.Environments,
		Variables:    suggestion.Variables,
		Rationale:    suggestion.Rationale,
	})
}
