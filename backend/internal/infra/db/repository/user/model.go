package userrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type User struct {
	model.Users
}

func (c *User) ToEntity() *entity.User {
	return &entity.User{
		ID:        c.ID,
		Email:     c.Email,
		Name:      c.Name,
		Role:      entity.UserRole(c.Role),
		CreatedAt: c.CreatedAt,
		LastLogin: misc.DerefTime(c.LastLogin),
		DeletedAt: misc.DerefTime(c.DeletedAt),
	}
}

type Users []User

func (us Users) ToEntities() *entity.Users {
	users := make(entity.Users, 0, len(us))
	for _, c := range us {
		user := c.ToEntity()
		if user == nil {
			continue
		}
		users = append(users, misc.GetValue(user))
	}

	return misc.ToPointer(users)
}
