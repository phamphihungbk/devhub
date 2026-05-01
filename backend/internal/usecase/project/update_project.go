package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type UpdateProjectInput struct {
	ID           string                    `json:"id" validate:"required,uuid"`
	Name         *string                   `json:"name" validate:"omitempty,min=0,max=100"`
	Description  *string                   `json:"description" validate:"omitempty,min=0,max=500"`
	OwnerTeamID  *string                   `json:"owner_team_id" validate:"omitempty,uuid"`
	Environments *[]CreateEnvironmentInput `json:"environments" validate:"omitempty,dive"`
}

func (u *projectUsecase) UpdateProject(ctx context.Context, input UpdateProjectInput) (project *entity.Project, err error) {
	const errLocation = "[usecase project/update_project UpdateProject] "
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

	var ownerTeamID *uuid.UUID
	if input.OwnerTeamID != nil {
		parsedOwnerTeamID, parseErr := uuid.Parse(*input.OwnerTeamID)
		if parseErr != nil {
			return nil, misc.WrapError(parseErr, errs.NewBadRequestError("invalid owner team id", nil))
		}
		ownerTeamID = &parsedOwnerTeamID
	}

	projectID := uuid.MustParse(input.ID)
	var updated *entity.Project
	if input.Name != nil || input.Description != nil || ownerTeamID != nil {
		updated, err = u.projectRepository.UpdateOne(ctx, repository.UpdateProjectInput{
			ID:          projectID,
			Name:        input.Name,
			Description: input.Description,
			OwnerTeamID: ownerTeamID,
		})
		if err != nil {
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update project", nil))
		}
	} else {
		updated, err = u.projectRepository.FindOne(ctx, projectID)
		if err != nil {
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find project", nil))
		}
	}

	if input.Environments != nil {
		if err := u.environmentRepository.DeleteByProjectID(ctx, updated.ID); err != nil {
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to delete project environments", nil))
		}

		for _, environmentInput := range *input.Environments {
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
				ProjectID:      updated.ID,
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
	}

	u.enrichProjectCreator(ctx, updated)
	if err := u.enrichProjectEnvironments(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
}
