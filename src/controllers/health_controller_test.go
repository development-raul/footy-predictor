package controllers

import (
	"github.com/development-raul/footy-predictor/src/utils/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthController_Check(t *testing.T) {
	// Initialization
	res:= httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := test.GetMockedContext(req, res)

	// Execution
	HealthController.Check(c)

	// Validation
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "I'm alive", res.Body.String())
}
