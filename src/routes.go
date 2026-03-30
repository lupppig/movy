package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lupppig/movy/internal/config"
	"github.com/lupppig/movy/internal/logger"
)




func Router(config config.BaseConfig, logger *logger.Logger) *gin.Engine  {
	r := gin.Default()

	r.GET(fmt.Sprintf("%s/health", config.API_VERSION), func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"version": "1.0.1",
			"healthy": true,
		})
	})

	return r
}