package handler

import (
	"devhub-backend/internal/domain/entity"
	releaseUsecase "devhub-backend/internal/usecase/release"
	"devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type findAllReleasesResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceID   string `json:"service_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PluginID    string `json:"plugin_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Tag         string `json:"tag" example:"v1.0.0"`
	Target      string `json:"target" example:"main"`
	Name        string `json:"name" example:"v1.0.0"`
	Notes       string `json:"notes" example:"First stable release."`
	HTMLURL     string `json:"html_url" example:"https://gitea.devhub.local/acme/service/releases/tag/v1.0.0"`
	ExternalRef string `json:"external_ref" example:"123"`
	Status      string `json:"status" example:"completed"`
	TriggeredBy string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary		List Releases
// @Description	List all releases for a service
// @Tags			Release
// @Produce		json
// @Success		200	{object}	httpresponse.SuccessResponse{data=[]findAllReleasesResponse,metadata=nil}	"List of releases"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/services/:service/releases [get]
func (h *releaseHandler) FindAllReleases(c *gin.Context) {
	serviceID := c.Param("service")

	releases, err := h.releaseUsecase.FindAllReleases(c.Request.Context(), releaseUsecase.FindAllReleasesInput{
		ServiceID: serviceID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindAllReleasesResponse(releases))
}

func (h *releaseHandler) newFindAllReleasesResponse(releases entity.Releases) []findAllReleasesResponse {
	if len(releases) == 0 {
		return []findAllReleasesResponse{}
	}

	response := make([]findAllReleasesResponse, 0, len(releases))
	for _, release := range releases {
		response = append(response, findAllReleasesResponse{
			ID:          release.ID.String(),
			ServiceID:   release.ServiceID.String(),
			PluginID:    release.PluginID.String(),
			Tag:         release.Tag,
			Target:      release.Target,
			Name:        release.Name,
			Notes:       release.Notes,
			HTMLURL:     release.HTMLURL,
			ExternalRef: release.ExternalRef,
			Status:      release.Status.String(),
			TriggeredBy: release.TriggeredBy.String(),
		})
	}

	return response
}
