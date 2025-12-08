package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updatePluginRequest struct {
	Name         *string   `json:"name" example:"Plugin Name" binding:"required"`
	Description  *string   `json:"description" example:"Plugin Description" binding:"required"`
	Environments *[]string `json:"environments" example:"[prod,dev,staging]" binding:"required"`
}

type updatePluginResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Plugin Name"`
	Description  string   `json:"description" example:"Plugin Description"`
	Environments []string `json:"environments" example:"[prod,dev,staging]"`
}

// @Summary		Update Plugin
// @Description	Update an existing plugin
// @Tags			Plugin
// @Accept			json
// @Produce		json
// @Param			request	body		updatePluginRequest													true	"Plugin update input"
// @Success		200		{object}	httpresponse.SuccessResponse{data=updatePluginResponse,metadata=nil}	    "Plugin updated"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/plugins/{plugin} [patch]
func (h *pluginHandler) UpdatePlugin(c *gin.Context) {
	pluginID := c.Param("plugin")
	var input updatePluginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	updatedPlugin, err := h.pluginUsecase.UpdatePlugin(c.Request.Context(), pluginUsecase.UpdatePluginInput{
		ID:           pluginID,
		Name:         input.Name,
		Description:  input.Description,
		Environments: input.Environments,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusOK, h.newUpdatePluginResponse(updatedPlugin))
}

func (h *pluginHandler) newUpdatePluginResponse(plugin *entity.Plugin) updatePluginResponse {
	if plugin == nil {
		return updatePluginResponse{}
	}

	environments := make([]string, len(plugin.Environments))
	for i, env := range plugin.Environments {
		environments[i] = string(env)
	}

	return updatePluginResponse{
		ID:           plugin.ID.String(),
		Name:         plugin.Name,
		Description:  plugin.Description,
		Environments: environments,
	}
}
