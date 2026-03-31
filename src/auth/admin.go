package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lupppig/movy/internal/openapi"
)

func (a *AuthDep) PromoteUserToAdmin(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, openapi.BadRequest{
			Code:    openapi.CodeInvalidInput,
			Message: "user id is required",
		})
		return
	}

	status, resp, uErr := a.PromoteUserService(userID)
	if uErr != nil {
		a.Logger.Error().Err(uErr).Msg(uErr.Error())
		c.JSON(status, uErr)
		return
	}
	c.JSON(status, resp)
}

func (a *AuthDep) GetUsers(c *gin.Context) {
	status, users, uErr := a.GetUsersService()
	if uErr != nil {
		a.Logger.Error().Err(uErr).Msg(uErr.Error())
		c.JSON(status, uErr)
		return
	}
	c.JSON(status, users)
}

func (a *AuthDep) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, openapi.BadRequest{
			Code:    openapi.CodeInvalidInput,
			Message: "user id is required",
		})
		return
	}

	status, user, uErr := a.GetUserService(userID)
	if uErr != nil {
		a.Logger.Error().Err(uErr).Msg(uErr.Error())
		c.JSON(status, uErr)
		return
	}
	c.JSON(status, user)
}

func (a *AuthDep) DemoteUserFromAdmin(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, openapi.BadRequest{
			Code:    openapi.CodeInvalidInput,
			Message: "user id is required",
		})
		return
	}

	status, resp, uErr := a.DemoteUserService(userID)
	if uErr != nil {
		a.Logger.Error().Err(uErr).Msg(uErr.Error())
		c.JSON(status, uErr)
		return
	}
	c.JSON(status, resp)
}
