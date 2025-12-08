package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/util/misc"
)

func (u *authUsecase) LogoutUser(ctx context.Context) (user *entity.User, err error) {
	const errLocation = "[usecase auth/logout_user LogoutUser] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	return nil, nil
}
