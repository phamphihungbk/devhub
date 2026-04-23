package usecase

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	repository "devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"
)

type FindAllApprovalRequestsInput struct {
	Limit     *int64                        `json:"limit" validate:"required,gte=1,lte=100"`
	Offset    *int64                        `json:"offset" validate:"required,gte=0"`
	SortBy    *string                       `json:"sort_by" validate:"required_with=SortOrder,omitempty,oneof=date status"`
	SortOrder *entity.SortOrder             `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	Status    *entity.ApprovalRequestStatus `json:"status" validate:"omitempty,oneof=pending approved rejected canceled"`
}

func (u *approvalUsecase) FindAllApprovalRequests(ctx context.Context, input FindAllApprovalRequestsInput) (requests entity.Page[entity.ApprovalRequest], err error) {
	const errLocation = "[usecase approval/find_all_requests FindAllApprovalRequests] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err := vInstance.Struct(input); err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	return entity.NewPage(u.findAllApprovalRequests(ctx, input))
}

func (u *approvalUsecase) findAllApprovalRequests(ctx context.Context, input FindAllApprovalRequestsInput) entity.PageProvider[entity.ApprovalRequest] {
	return func() ([]entity.ApprovalRequest, entity.PageProvider[entity.ApprovalRequest], entity.Pagination, error) {
		requests, count, err := u.approvalRepository.FindAllApprovalRequests(ctx, repository.FindAllApprovalRequestsFilter{
			Limit:     input.Limit,
			Offset:    input.Offset,
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
			Status:    input.Status,
		})
		if err != nil {
			return entity.ApprovalRequests{}, nil, entity.Pagination{}, errs.NewInternalServerError("failed to fetch approval requests", nil)
		}
		if requests == nil || len(misc.GetValue(requests)) == 0 {
			return entity.ApprovalRequests{}, nil, entity.Pagination{}, nil
		}

		pagination := entity.NewPagination(count, misc.GetValue(input.Limit), misc.GetValue(input.Offset))
		nextSearchCriteria := FindAllApprovalRequestsInput{
			Limit:     input.Limit,
			Offset:    misc.ToPointer((*input.Limit) + (*input.Offset)),
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
			Status:    input.Status,
		}
		return misc.GetValue(requests), u.findAllApprovalRequests(ctx, nextSearchCriteria), pagination, nil
	}
}
