package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateProjectInput struct {
	Name         string                   `json:"name" validate:"required,min=2,max=100"`
	Description  *string                  `json:"description" validate:"omitempty,min=0,max=500"`
	OwnerTeamID  string                   `json:"owner_team_id" validate:"required,uuid"`
	Environments []CreateEnvironmentInput `json:"environments" validate:"omitempty,dive"`
	CreatedBy    string                   `json:"created_by" validate:"required,uuid"`
}

type CreateEnvironmentConfigInput struct {
	Domain       string `json:"domain" validate:"required"`
	APIDomain    string `json:"api_domain" validate:"required"`
	Region       string `json:"region" validate:"required"`
	IngressClass string `json:"ingress_class" validate:"required"`
	PublicAccess bool   `json:"public_access"`
}

type CreateEnvironmentInput struct {
	Name           string                        `json:"name" validate:"required,min=1,max=255"`
	Tier           string                        `json:"tier" validate:"required,oneof=development staging production"`
	Cluster        string                        `json:"cluster" validate:"required,min=1,max=255"`
	Namespace      string                        `json:"namespace" validate:"required,min=1,max=255"`
	ArgoCDInstance string                        `json:"argocd_instance" validate:"required,min=1,max=255"`
	Config         *CreateEnvironmentConfigInput `json:"config" validate:"omitempty"`
}

func (u *projectUsecase) CreateProject(ctx context.Context, input CreateProjectInput) (project *entity.Project, err error) {
	const errLocation = "[usecase project/create_project CreateProject] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	ownerTeamID, err := uuid.Parse(input.OwnerTeamID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid owner team id", nil))
	}

	project = &entity.Project{
		Name:        input.Name,
		Description: misc.GetValue(input.Description),
		OwnerTeamID: ownerTeamID,
		CreatedBy:   uuid.MustParse(input.CreatedBy),
	}

	created, err := u.projectRepository.CreateOne(ctx, project)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create project", nil))
	}

	for _, environmentInput := range input.Environments {
		tier, err := new(entity.EnvironmentTier).Parse(environmentInput.Tier)
		if err != nil {
			return nil, misc.WrapError(err, errs.NewBadRequestError("invalid environment tier", map[string]string{"details": err.Error()}))
		}

		config := entity.EnvironmentConfig{}
		if environmentInput.Config != nil {
			config = entity.EnvironmentConfig{
				Domain:       environmentInput.Config.Domain,
				APIDomain:    environmentInput.Config.APIDomain,
				Region:       environmentInput.Config.Region,
				IngressClass: environmentInput.Config.IngressClass,
				PublicAccess: environmentInput.Config.PublicAccess,
			}
		}

		_, err = u.environmentRepository.CreateOne(ctx, &entity.Environment{
			ProjectID:      created.ID,
			Name:           environmentInput.Name,
			Tier:           tier,
			Cluster:        environmentInput.Cluster,
			Namespace:      environmentInput.Namespace,
			ArgocdInstance: environmentInput.ArgoCDInstance,
			Config:         config,
		})
		if err != nil {
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create project environment", nil))
		}
	}

	return created, nil
}
