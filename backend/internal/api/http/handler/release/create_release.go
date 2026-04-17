package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	releaseUsecase "devhub-backend/internal/usecase/release"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createReleaseRequest struct {
	PluginID string `json:"plugin_id" example:"123e4567-e89b-12d3-a456-426614174000" binding:"required"`
	Tag    string `json:"tag" example:"v1.0.0" binding:"required"`
	Target string `json:"target" example:"main"`
	Name   string `json:"name" example:"v1.0.0"`
	Notes  string `json:"notes" example:"First stable release."`
}

type createReleaseResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ProjectID   string `json:"project_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PluginID    string `json:"plugin_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Tag         string `json:"tag" example:"v1.0.0"`
	Target      string `json:"target" example:"main"`
	Name        string `json:"name" example:"v1.0.0"`
	Notes       string `json:"notes" example:"First stable release."`
	HTMLURL     string `json:"html_url" example:"https://gitea.devhub.local/acme/service/releases/tag/v1.0.0"`
	ExternalRef string `json:"external_ref" example:"123"`
	TriggeredBy string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary		Create Release
// @Description	Create a Git tag and Gitea release for a project repository
// @Tags			Release
// @Accept			json
// @Produce		json
// @Param			request	body		createReleaseRequest													true	"Release creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createReleaseResponse,metadata=nil}	"Release created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		409		{object}	httpresponse.ErrorResponse{data=nil}									"Conflict"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects/:project/releases [post]
func (h *releaseHandler) CreateRelease(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	projectID := c.Param("project")
	var input createReleaseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	release, err := h.releaseUsecase.CreateRelease(c.Request.Context(), releaseUsecase.CreateReleaseInput{
		ProjectID:   projectID,
		PluginID:    input.PluginID,
		Tag:         input.Tag,
		Target:      input.Target,
		Name:        input.Name,
		Notes:       input.Notes,
		TriggeredBy: userID.(string),
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateReleaseResponse(release))
}

func (h *releaseHandler) newCreateReleaseResponse(release *entity.Release) createReleaseResponse {
	if release == nil {
		return createReleaseResponse{}
	}

	return createReleaseResponse{
		ID:          release.ID.String(),
		ProjectID:   release.ProjectID.String(),
		PluginID:    release.PluginID.String(),
		Tag:         release.Tag,
		Target:      release.Target,
		Name:        release.Name,
		Notes:       release.Notes,
		HTMLURL:     release.HTMLURL,
		ExternalRef: release.ExternalRef,
		TriggeredBy: release.TriggeredBy.String(),
	}
}
