package controllers

import (
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/services"
	"github.com/development-raul/footy-predictor/src/utils"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCountryService struct {
	FuncCreate func(req *countries.CreateCountryInput) resterror.RestErrorI
}

func (m MockCountryService) Create(req *countries.CreateCountryInput) resterror.RestErrorI {
	return m.FuncCreate(req)
}

func TestCountryController_Create(t *testing.T) {
	testCases := []struct {
		title string
		reqBody io.Reader
		serviceMock services.CountryServiceI
		expectedStatus int
		expectedRes string
	}{
		{
			title: "error required name",
			reqBody: strings.NewReader(`{}`),
			serviceMock: nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes: `{"error":{"name":["The name field is required."]},"code":400}`,
		},
		{
			title: "error invalid active",
			reqBody: strings.NewReader(`{"name":"England","active":2}`),
			serviceMock: nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes: `{"error":{"active":["The field: 'active' must be one of [0 1]"]},"code":400}`,
		},
		{
			title: "error CountryService.Create",
			reqBody: strings.NewReader(`{"name":"England","active":1}`),
			serviceMock: &MockCountryService{
				FuncCreate: func(req *countries.CreateCountryInput) resterror.RestErrorI {
					return resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes: `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title: "success",
			reqBody: strings.NewReader(`{"name":"England","active":1}`),
			serviceMock: &MockCountryService{
				FuncCreate: func(req *countries.CreateCountryInput) resterror.RestErrorI {
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedRes: `{"message":"SUCCESS","code":201}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "https://localhost:8000/v1/countries", testCase.reqBody)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			c := utils.GetMockedContext(req, res)

			services.CountryService = testCase.serviceMock
			CountryController.Create(c)

			assert.Equal(t, testCase.expectedStatus, res.Code)
			assert.Equal(t, testCase.expectedRes, res.Body.String())
		})
	}
}