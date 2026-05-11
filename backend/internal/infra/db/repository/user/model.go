package userrepo

import (
	"encoding/json"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type User struct {
	model.Users
	RolesJSON *string `db:"roles_json"`
}

func (c *User) ToEntity() *entity.User {
	roles := []string{}
	if c.RolesJSON != nil {
		err := json.Unmarshal([]byte(misc.GetValue(c.RolesJSON)), &roles)
		if err != nil {
			return nil
		}
	}

	return &entity.User{
		ID:           c.ID,
		Name:         c.Name,
		Email:        c.Email,
		PasswordHash: c.PasswordHash,
		Roles:        roles,
		TeamID:       c.TeamID,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		DeletedAt:    misc.DerefTime(c.DeletedAt),
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

type Role struct {
	model.Roles
}

func (r *Role) ToEntity() *entity.Role {
	return &entity.Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: misc.GetValue(r.Description),
		CreatedAt:   r.CreatedAt,
	}
}

type Roles []Role

func (rs Roles) ToEntities() []entity.Role {
	roles := make([]entity.Role, 0, len(rs))
	for _, item := range rs {
		role := item.ToEntity()
		if role == nil {
			continue
		}
		roles = append(roles, misc.GetValue(role))
	}

	return roles
}
