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

type createPluginRequest struct {
	Name        string `json:"name" example:"Plugin Name" binding:"required"`
	Version     string `json:"version" example:"1.0.0" binding:"required"`
	Type        string `json:"type" example:"scaffolder" binding:"required"`
	Entrypoint  string `json:"entrypoint" example:"/app/plugins/scaffolders/go_http_api/action.py" binding:"required"`
	Scope       string `json:"scope" example:"global" binding:"required"`
	Description string `json:"description" example:"Plugin Description" binding:"required"`
}

type createPluginResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string `json:"name" example:"Plugin Name"`
	Version     string `json:"version" example:"1.0.0"`
	Type        string `json:"type" example:"scaffolder"`
	Entrypoint  string `json:"entrypoint" example:"/app/plugins/scaffolders/go_http_api/action.py"`
	Scope       string `json:"scope" example:"global"`
	Description string `json:"description" example:"Plugin Description"`
}

// @Summary		Create Plugin
// @Description	Create a new plugin
// @Tags			Plugin
// @Accept			json
// @Produce		json
// @Param			request	body		createPluginRequest													true	"Plugin creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createPluginResponse,metadata=nil}	"Plugin created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/plugins [post]
func (h *pluginHandler) CreatePlugin(c *gin.Context) {
	var input createPluginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdPlugin, err := h.pluginUsecase.CreatePlugin(c.Request.Context(), pluginUsecase.CreatePluginInput{
		Name:        input.Name,
		Version:     input.Version,
		Type:        input.Type,
		Entrypoint:  input.Entrypoint,
		Scope:       input.Scope,
		Description: input.Description,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreatePluginResponse(createdPlugin))
}

func (h *pluginHandler) newCreatePluginResponse(plugin *entity.Plugin) createPluginResponse {
	if plugin == nil {
		return createPluginResponse{}
	}

	return createPluginResponse{
		ID:          plugin.ID.String(),
		Version:     plugin.Version,
		Name:        plugin.Name,
		Type:        plugin.Type.String(),
		Entrypoint:  plugin.Entrypoint,
		Scope:       plugin.Scope,
		Description: plugin.Description,
	}
}
