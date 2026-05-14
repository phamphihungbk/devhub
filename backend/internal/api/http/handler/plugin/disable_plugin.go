package handler

import (
	"devhub-backend/internal/domain/entity"
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	"devhub-backend/internal/util/httpresponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type disablePluginResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string `json:"name" example:"Plugin Name"`
	Type        string `json:"type" example:"scaffolder"`
	Version     string `json:"version" example:"1.0.0"`
	Runtime     string `json:"runtime" example:"python"`
	Entrypoint  string `json:"entrypoint" example:"/app/plugins/scaffolders/go_http_api/action.py"`
	Enabled     bool   `json:"enabled" example:"true"`
	Description string `json:"description" example:"Plugin Description"`
}

// @Summary		Disable Plugin
// @Description	Disable an existing plugin
// @Tags			Plugin
// @Produce		json
// @Success		200		{object}	httpresponse.SuccessResponse{data=disablePluginResponse,metadata=nil}	"Plugin disabled"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404		{object}	httpresponse.ErrorResponse{data=nil}									"Plugin not found"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/plugins/{plugin}/disable [patch]
func (h *pluginHandler) DisablePlugin(c *gin.Context) {
	pluginID := c.Param("plugin")
	enabled := false

	updatedPlugin, err := h.pluginUsecase.UpdatePlugin(c.Request.Context(), pluginUsecase.UpdatePluginInput{
		ID:      pluginID,
		Enabled: &enabled,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusOK, h.newDisablePluginResponse(updatedPlugin))
}

func (h *pluginHandler) newDisablePluginResponse(plugin *entity.Plugin) disablePluginResponse {
	if plugin == nil {
		return disablePluginResponse{}
	}

	return disablePluginResponse{
		ID:          plugin.ID.String(),
		Name:        plugin.Name,
		Description: plugin.Description,
		Version:     plugin.Version,
		Type:        plugin.Type.String(),
		Runtime:     plugin.Runtime.String(),
		Entrypoint:  plugin.Entrypoint,
		Enabled:     plugin.Enabled,
	}
}
