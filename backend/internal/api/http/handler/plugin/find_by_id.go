package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/util/httpresponse"

	pluginUsecase "devhub-backend/internal/usecase/plugin"

	"github.com/gin-gonic/gin"
)

type findOnePluginResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Plugin Name"`
	Description  string   `json:"description" example:"Plugin Description"`
	Environments []string `json:"environments" example:"[prod,dev,staging]"`
}

// @Summary		Find Plugin by ID
// @Description	Retrieve plugin details by its ID
// @Tags			Plugin
// @Produce		json
// @Param			plugin	path		string																	true	"Plugin ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=findOnePluginResponse,metadata=nil}	"Plugin found"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Plugin not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/plugins/{plugin} [get]
func (h *pluginHandler) FindPluginByID(c *gin.Context) {
	pluginID := c.Param("plugin")
	plugin, err := h.pluginUsecase.FindOnePlugin(c.Request.Context(), pluginUsecase.FindOnePluginInput{
		ID: pluginID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindOnePluginResponse(plugin))
}

func (h *pluginHandler) newFindOnePluginResponse(plugin *entity.Plugin) findOnePluginResponse {
	if plugin == nil {
		return findOnePluginResponse{}
	}

	environments := make([]string, len(plugin.Environments))
	for i, env := range plugin.Environments {
		environments[i] = string(env)
	}

	return findOnePluginResponse{
		ID:           plugin.ID.String(),
		Name:         plugin.Name,
		Description:  plugin.Description,
		Environments: environments,
	}
}
