package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lupppig/movy/internal/openapi"
	"github.com/lupppig/movy/internal/role"
)

func IDORMiddleware(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, openapi.UnauthorizedError{
				Code:    openapi.CodeUnauthorized,
				Message: "authentication required",
			})
			c.Abort()
			return
		}

		userRole, _ := c.Get("role")
		if userRole == role.Admin {
			c.Next()
			return
		}

		resourceID := c.Param(paramKey)
		if resourceID == "" {
			c.JSON(http.StatusBadRequest, openapi.BadRequestError{
				Code:    openapi.CodeInvalidInput,
				Message: "resource id is required",
			})
			c.Abort()
			return
		}

		if userID != resourceID {
			c.JSON(http.StatusForbidden, openapi.ForbiddenError{
				Code:    openapi.CodeUnauthorized,
				Message: "you do not have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
