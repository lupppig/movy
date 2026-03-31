package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lupppig/movy/internal/logger"
	"github.com/lupppig/movy/internal/openapi"
	"github.com/lupppig/movy/internal/utils"
)

func AuthMiddleware(secret string, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn().Msg("missing authorization header")
			c.JSON(http.StatusUnauthorized, openapi.BadRequest{
				Code:    openapi.CodeUnauthorized,
				Message: "missing authorization header",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logger.Warn().Msg("invalid authorization header format")
			c.JSON(http.StatusUnauthorized, openapi.BadRequest{
				Code:    openapi.CodeUnauthorized,
				Message: "invalid authorization header format, expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString, secret)
		if err != nil {
			logger.Warn().Err(err).Msg("invalid or expired token")
			c.JSON(http.StatusUnauthorized, openapi.BadRequest{
				Code:    openapi.CodeUnauthorized,
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, openapi.Error{
				Code:    openapi.CodeUnauthorized,
				Message: "access denied",
			})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, openapi.Error{
			Code:    openapi.CodeUnauthorized,
			Message: "insufficient permissions",
		})
		c.Abort()
	}
}
