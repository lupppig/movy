package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lupppig/movy/auth"
	"github.com/lupppig/movy/internal/config"
	"github.com/lupppig/movy/internal/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)




func Router(config config.BaseConfig, logger *logger.Logger, db *sql.DB) *gin.Engine  {
	r := gin.Default()

	r.GET(fmt.Sprintf("%s/health", config.API_VERSION), func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"version": "1.0.1",
			"status": "alive",
		})
	})

	r.StaticFile("/internal/openapi/openapi.yaml", "./internal/openapi/openapi.yaml")

	// swagger documentation...
	url := ginSwagger.URL("/internal/openapi/openapi.yaml") // Points to the static route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))



	// User auth routes
	api := r.Group(config.API_VERSION)
	{
		authentication:= api.Group("/auth") 
		{
			a := auth.AuthDep{Logger: logger, Config: &config, DB: db}
			authentication.POST("/signup", a.RegisterUser)
		}
	}
	return r
}