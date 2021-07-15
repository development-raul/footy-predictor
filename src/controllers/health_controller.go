package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthControllerInterface interface {
	Check(c *gin.Context)
}

type healthController struct{}

var HealthController healthControllerInterface = &healthController{}

// Check godoc
// @Summary Health check endpoint.
// @Description Will return a 200 status code if the application is up and running
// @ID health-check
// @Tags Health Check
// @Success 200 {string} string
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Router / [get]
func (hc *healthController) Check(c *gin.Context) {
	c.String(http.StatusOK, "I'm alive")
}