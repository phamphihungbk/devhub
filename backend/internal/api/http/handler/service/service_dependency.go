package handler

import (
	"net/http"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	serviceUsecase "devhub-backend/internal/usecase/service"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type serviceDependencyRequest struct {
	DependsOnServiceID string         `json:"depends_on_service_id" binding:"required"`
	Type               string         `json:"type" binding:"required"`
	Protocol           string         `json:"protocol"`
	Port               *int           `json:"port"`
	Path               string         `json:"path"`
	Config             map[string]any `json:"config"`
}

type serviceDependencyResponse struct {
	ID                 string                   `json:"id"`
	ServiceID          string                   `json:"service_id"`
	DependsOnServiceID string                   `json:"depends_on_service_id"`
	DependsOnService   *findAllServicesResponse `json:"depends_on_service,omitempty"`
	Type               string                   `json:"type"`
	Protocol           string                   `json:"protocol"`
	Port               *int                     `json:"port"`
	Path               string                   `json:"path"`
	Config             map[string]any           `json:"config"`
	CreatedBy          string                   `json:"created_by"`
}

// @Summary		List Service Dependencies
// @Description	List services that the selected service depends on
// @Tags			Service
// @Produce		json
// @Success		200	{object}	httpresponse.SuccessResponse{data=[]serviceDependencyResponse,metadata=nil}	"List of service dependencies"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}										"Bad request"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}										"Internal server error"
// @Router			/services/:service/dependencies [get]
func (h *serviceHandler) FindServiceDependencies(c *gin.Context) {
	dependencies, err := h.serviceUsecase.FindServiceDependencies(c.Request.Context(), serviceUsecase.FindServiceDependenciesInput{
		ServiceID: c.Param("service"),
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newServiceDependencyResponses(dependencies))
}

// @Summary		Create Service Dependency
// @Description	Wire one service to another service in the same project
// @Tags			Service
// @Accept			json
// @Produce		json
// @Param			request	body		serviceDependencyRequest													true	"Service dependency input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=serviceDependencyResponse,metadata=nil}	"Service dependency created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/services/:service/dependencies [post]
func (h *serviceHandler) CreateServiceDependency(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	var input serviceDependencyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		httpresponse.Error(c, misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()})))
		return
	}

	dependency, err := h.serviceUsecase.CreateServiceDependency(c.Request.Context(), serviceUsecase.CreateServiceDependencyInput{
		ServiceID:          c.Param("service"),
		DependsOnServiceID: input.DependsOnServiceID,
		Type:               input.Type,
		Protocol:           input.Protocol,
		Port:               input.Port,
		Path:               input.Path,
		Config:             input.Config,
		CreatedBy:          userID.(string),
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newServiceDependencyResponse(dependency))
}

// @Summary		Delete Service Dependency
// @Description	Remove a service-to-service wiring entry
// @Tags			Service
// @Produce		json
// @Success		200	{object}	httpresponse.SuccessResponse{data=serviceDependencyResponse,metadata=nil}	"Service dependency deleted"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/services/:service/dependencies/:dependency [delete]
func (h *serviceHandler) DeleteServiceDependency(c *gin.Context) {
	dependency, err := h.serviceUsecase.DeleteServiceDependency(c.Request.Context(), serviceUsecase.DeleteServiceDependencyInput{
		ServiceID:    c.Param("service"),
		DependencyID: c.Param("dependency"),
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newServiceDependencyResponse(dependency))
}

func (h *serviceHandler) newServiceDependencyResponses(dependencies entity.ServiceDependencies) []serviceDependencyResponse {
	if len(dependencies) == 0 {
		return []serviceDependencyResponse{}
	}

	response := make([]serviceDependencyResponse, 0, len(dependencies))
	for _, dependency := range dependencies {
		response = append(response, h.newServiceDependencyResponse(&dependency))
	}

	return response
}

func (h *serviceHandler) newServiceDependencyResponse(dependency *entity.ServiceDependency) serviceDependencyResponse {
	if dependency == nil {
		return serviceDependencyResponse{}
	}

	var dependsOn *findAllServicesResponse
	if dependency.DependsOnService != nil {
		dependsOn = &findAllServicesResponse{
			ID:        dependency.DependsOnService.ID.String(),
			ProjectID: dependency.DependsOnService.ProjectID.String(),
			Name:      dependency.DependsOnService.Name,
			RepoURL:   dependency.DependsOnService.RepoURL,
		}
	}

	return serviceDependencyResponse{
		ID:                 dependency.ID.String(),
		ServiceID:          dependency.ServiceID.String(),
		DependsOnServiceID: dependency.DependsOnServiceID.String(),
		DependsOnService:   dependsOn,
		Type:               dependency.Type,
		Protocol:           dependency.Protocol,
		Port:               dependency.Port,
		Path:               dependency.Path,
		Config:             dependency.Config,
		CreatedBy:          dependency.CreatedBy.String(),
	}
}
