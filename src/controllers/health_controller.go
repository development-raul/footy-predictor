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

func (hc *healthController) Check(c *gin.Context) {
	c.String(http.StatusOK, "I'm alive")
}