package handler

import "github.com/gin-gonic/gin"

type livenessResponse struct {
	Status string `json:"status" example:"OK"`
}

func (h *healthCheckHandler) Liveness(c *gin.Context) {
	httpresponse.Success(c, h.newLivenessResponse())
}

func (h *healthCheckHandler) newLivenessResponse() livenessResponse {
	return livenessResponse{
		Status: "OK",
	}
}
