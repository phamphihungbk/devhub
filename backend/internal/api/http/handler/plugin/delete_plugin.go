package handler

import (
	"devhub-backend/internal/util/httpresponse"

	pluginUsecase "devhub-backend/internal/usecase/plugin"

	"github.com/gin-gonic/gin"
)

// @Summary		Delete Plugin by ID
// @Description	Delete a plugin by its ID
// @Tags			Plugin
// @Produce		json
// @Param			plugin	path		string																	true	"Plugin ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=nil,metadata=nil}	"Plugin deleted"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Plugin not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/plugins/{plugin} [delete]
func (h *pluginHandler) DeletePlugin(c *gin.Context) {
	pluginID := c.Param("plugin")
	_, err := h.pluginUsecase.DeletePlugin(c.Request.Context(), pluginUsecase.DeletePluginInput{
		ID: pluginID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, nil)
}
