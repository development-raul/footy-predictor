package test

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func GetMockedContext(req *http.Request, res *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(res)
	c.Request = req
	return c
}