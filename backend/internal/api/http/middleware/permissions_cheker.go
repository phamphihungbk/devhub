package middleware

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (m *middleware) RequirePermissions(permissions ...entity.PermissionName) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(permissions) == 0 {
			ctx.Next()
			return
		}

		userIDValue, exists := ctx.Get("user_id")
		if !exists {
			httpresponse.Error(ctx, errs.NewUnauthorizedError("Missing authenticated user", nil))
			return
		}

		userIDString, ok := userIDValue.(string)
		if !ok {
			httpresponse.Error(ctx, errs.NewForbiddenError("Invalid user context", nil))
			return
		}

		userID, err := uuid.Parse(userIDString)
		if err != nil {
			httpresponse.Error(ctx, errs.NewForbiddenError("Invalid user context", map[string]string{"user_id": userIDString}))
			return
		}

		// TODO: store permission on redis cache so no need to query again
		allowed, err := m.userRepository.HasPermissions(ctx.Request.Context(), userID, permissions)
		if err != nil {
			httpresponse.Error(ctx, errs.NewInternalServerError("Failed to check permissions", nil))
			return
		}
		if !allowed {
			httpresponse.Error(ctx, errs.NewForbiddenError("Insufficient permissions", map[string]any{
				"permissions": permissions,
			}))
			return
		}

		ctx.Next()
	}
}
