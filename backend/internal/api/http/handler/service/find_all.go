package handler

import (
	"devhub-backend/internal/domain/entity"
	serviceUsecase "devhub-backend/internal/usecase/service"
	"devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type findAllServicesResponse struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ProjectID string `json:"project_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string `json:"name" example:"payment"`
	RepoURL   string `json:"repo_url" example:"https://gitea.devhub.local/platform/payment.git"`
}

// @Summary		List Services
// @Description	List all services for a project
// @Tags			Service
// @Produce		json
// @Success		200	{object}	httpresponse.SuccessResponse{data=[]findAllServicesResponse,metadata=nil}	"List of services"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects/:project/services [get]
func (h *serviceHandler) FindAllServices(c *gin.Context) {
	projectID := c.Param("project")

	services, err := h.serviceUsecase.FindAllServices(c.Request.Context(), serviceUsecase.FindAllServicesInput{
		ProjectID: projectID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindAllServicesResponse(services))
}

func (h *serviceHandler) newFindAllServicesResponse(services entity.Services) []findAllServicesResponse {
	if len(services) == 0 {
		return []findAllServicesResponse{}
	}

	response := make([]findAllServicesResponse, 0, len(services))
	for _, service := range services {
		response = append(response, findAllServicesResponse{
			ID:        service.ID.String(),
			ProjectID: service.ProjectID.String(),
			Name:      service.Name,
			RepoURL:   service.RepoURL,
		})
	}

	return response
}
