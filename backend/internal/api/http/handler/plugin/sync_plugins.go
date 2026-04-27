package handler

import (
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	httpresponse "devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

const defaultPluginsDir = "../plugins"

type syncPluginsResponse struct {
	Discovered int `json:"discovered"`
	Created    int `json:"created"`
	Updated    int `json:"updated"`
}

func (h *pluginHandler) SyncPlugins(c *gin.Context) {
	output, err := h.pluginUsecase.SyncRegistry(c.Request.Context(), pluginUsecase.SyncRegistryInput{
		PluginsDir: defaultPluginsDir,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, syncPluginsResponse{
		Discovered: output.Discovered,
		Created:    output.Created,
		Updated:    output.Updated,
	})
}
