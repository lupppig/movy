package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lupppig/movy/internal/config"
	"github.com/lupppig/movy/internal/logger"
	"github.com/lupppig/movy/internal/openapi"
)

type AuthDep struct {
	Logger *logger.Logger
	Config *config.BaseConfig
	DB     *sql.DB
}

func (a *AuthDep) RegisterUser(c *gin.Context) {
	var req = &openapi.SignupRequest{}
	if err := c.BindJSON(req); err != nil {
		a.Logger.Error().Err(err).Msg("failed to bind user json request to userReq struct variable")
		c.JSON(http.StatusBadRequest, openapi.BadRequest{
			Code:    openapi.CodeInvalidInput,
			Message: "failed to decode request body",
		})
	}

	status, resp, uErr := a.CreateUserService(*req)
	if uErr != nil {
		a.Logger.Error().Err(uErr).Msg(uErr.Error())
		c.JSON(status, uErr)
		return
	}
	c.JSON(status, resp)
}

func (a *AuthDep) SignInUser(c *gin.Context) {
	var req = &openapi.SigninRequest{}
	if err := c.BindJSON(req); err != nil {
		a.Logger.Error().Err(err).Msg("failed to bind signin request")
		c.JSON(http.StatusBadRequest, openapi.BadRequest{
			Code:    openapi.CodeInvalidInput,
			Message: "failed to decode request body",
		})
		return
	}

	status, resp, uErr := a.SignInUserService(*req)
	if uErr != nil {
		a.Logger.Error().Err(uErr).Msg(uErr.Error())
		c.JSON(status, uErr)
		return
	}
	c.JSON(status, resp)
}
