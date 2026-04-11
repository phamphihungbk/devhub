package refreshtokenrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type RefreshToken struct {
	model.RefreshTokens
}

func (c *RefreshToken) ToEntity() *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        c.ID,
		UserID:    c.UserID,
		Token:     c.Token,
		CreatedAt: c.CreatedAt,
		ExpiresAt: c.ExpiresAt,
		DeletedAt: misc.DerefTime(c.DeletedAt),
	}
}
