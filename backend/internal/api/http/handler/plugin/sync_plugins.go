package handler

import (
	"devhub-backend/internal/domain/errs"
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	httpresponse "devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

const defaultPluginsDir = "../plugins"

type syncPluginsRequest struct {
	PluginsDir string `json:"plugins_dir" example:"../plugins"`
}

type syncPluginsResponse struct {
	Discovered int `json:"discovered"`
	Created    int `json:"created"`
	Updated    int `json:"updated"`
}

// @Summary		Sync Plugins
// @Description	Discover plugin manifests and create or update plugin records
// @Tags			Plugin
// @Accept			json
// @Produce		json
// @Param			request	body		syncPluginsRequest													false	"Plugin sync input"
// @Success		200		{object}	httpresponse.SuccessResponse{data=syncPluginsResponse,metadata=nil}	"Plugins synced"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}								"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}								"Internal server error"
// @Router			/plugins/sync [post]
func (h *pluginHandler) SyncPlugins(c *gin.Context) {
	input := syncPluginsRequest{PluginsDir: defaultPluginsDir}
	if c.Request.Body != nil && c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&input); err != nil {
			err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
			httpresponse.Error(c, err)
			return
		}
	}

	output, err := h.pluginUsecase.SyncRegistry(c.Request.Context(), pluginUsecase.SyncRegistryInput{
		PluginsDir: input.PluginsDir,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusOK, syncPluginsResponse{
		Discovered: output.Discovered,
		Created:    output.Created,
		Updated:    output.Updated,
	})
}
