package middleware

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

func (m *middleware) RequirePermissions(permissions ...entity.Permission) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleValue, exists := ctx.Get("role")
		if !exists {
			httpresponse.Error(ctx, errs.NewUnauthorizedError("Missing authenticated role", nil))
			return
		}

		roleString, ok := roleValue.(string)
		if !ok {
			httpresponse.Error(ctx, errs.NewForbiddenError("Invalid role context", nil))
			return
		}

		role, err := new(entity.UserRole).Parse(roleString)
		if err != nil {
			httpresponse.Error(ctx, errs.NewForbiddenError("Invalid user role", map[string]string{"role": roleString}))
			return
		}

		if !role.HasPermissions(permissions...) {
			httpresponse.Error(ctx, errs.NewForbiddenError("Insufficient permissions", map[string]any{
				"role":        role.String(),
				"permissions": permissions,
			}))
			return
		}

		ctx.Next()
	}
}
