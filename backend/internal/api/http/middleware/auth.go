package middleware

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func (m *middleware) Auth(tokenSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
			httpresponse.Error(ctx, errs.NewUnauthorizedError("Missing or malformed token", nil))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &entity.AccessTokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})

		if err != nil || !token.Valid {
			httpresponse.Error(ctx, misc.WrapError(err, errs.NewUnauthorizedError("Invalid or expired token", nil)))
			return
		}

		ctx.Set("user_id", claims.Subject)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
