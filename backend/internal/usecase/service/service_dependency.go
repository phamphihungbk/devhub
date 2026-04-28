package usecase

import (
	"context"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type FindServiceDependenciesInput struct {
	ServiceID string `json:"service_id" validate:"required,uuid"`
}

type CreateServiceDependencyInput struct {
	ServiceID          string         `json:"service_id" validate:"required,uuid"`
	DependsOnServiceID string         `json:"depends_on_service_id" validate:"required,uuid"`
	Type               string         `json:"type" validate:"required,oneof=http grpc queue database"`
	Protocol           string         `json:"protocol" validate:"omitempty,oneof=http https grpc tcp udp"`
	Port               *int           `json:"port" validate:"omitempty,min=1,max=65535"`
	Path               string         `json:"path" validate:"omitempty"`
	Config             map[string]any `json:"config"`
	CreatedBy          string         `json:"created_by" validate:"required,uuid"`
}

type DeleteServiceDependencyInput struct {
	ServiceID    string `json:"service_id" validate:"required,uuid"`
	DependencyID string `json:"dependency_id" validate:"required,uuid"`
}

func (u *serviceUsecase) FindServiceDependencies(ctx context.Context, input FindServiceDependenciesInput) (entity.ServiceDependencies, error) {
	if err := validateServiceDependencyInput(input); err != nil {
		return nil, err
	}

	serviceID := uuid.MustParse(input.ServiceID)
	dependencies, err := u.serviceRepository.FindDependencies(ctx, serviceID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to fetch service dependencies", nil))
	}
	if dependencies == nil {
		return entity.ServiceDependencies{}, nil
	}

	return misc.GetValue(dependencies), nil
}

func (u *serviceUsecase) CreateServiceDependency(ctx context.Context, input CreateServiceDependencyInput) (*entity.ServiceDependency, error) {
	if err := validateServiceDependencyInput(input); err != nil {
		return nil, err
	}

	serviceID := uuid.MustParse(input.ServiceID)
	dependsOnServiceID := uuid.MustParse(input.DependsOnServiceID)
	if serviceID == dependsOnServiceID {
		return nil, errs.NewBadRequestError("a service cannot depend on itself", nil)
	}

	service, err := u.serviceRepository.FindOne(ctx, serviceID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid service", nil))
	}

	dependsOnService, err := u.serviceRepository.FindOne(ctx, dependsOnServiceID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid dependent service", nil))
	}

	if service.ProjectID != dependsOnService.ProjectID {
		return nil, errs.NewBadRequestError("service dependencies must stay within the same project", nil)
	}

	config := input.Config
	if config == nil {
		config = map[string]any{}
	}

	dependency, err := u.serviceRepository.CreateDependency(ctx, &entity.ServiceDependency{
		ServiceID:          serviceID,
		DependsOnServiceID: dependsOnServiceID,
		Type:               strings.TrimSpace(input.Type),
		Protocol:           strings.TrimSpace(input.Protocol),
		Port:               input.Port,
		Path:               strings.TrimSpace(input.Path),
		Config:             config,
		CreatedBy:          uuid.MustParse(input.CreatedBy),
	})
	if err != nil {
		return nil, err
	}

	dependency.DependsOnService = dependsOnService
	return dependency, nil
}

func (u *serviceUsecase) DeleteServiceDependency(ctx context.Context, input DeleteServiceDependencyInput) (*entity.ServiceDependency, error) {
	if err := validateServiceDependencyInput(input); err != nil {
		return nil, err
	}

	dependency, err := u.serviceRepository.DeleteDependency(ctx, uuid.MustParse(input.ServiceID), uuid.MustParse(input.DependencyID))
	if err != nil {
		return nil, err
	}

	return dependency, nil
}

func validateServiceDependencyInput(input any) error {
	vInstance, err := validator.NewValidator(validator.WithTagNameFunc(validator.JSONTagNameFunc))
	if err != nil {
		return misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err := vInstance.Struct(input); err != nil {
		return misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	return nil
}
