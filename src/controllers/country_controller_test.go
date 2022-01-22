package controllers

import (
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/services"
	"github.com/development-raul/footy-predictor/src/utils"
	"github.com/development-raul/footy-predictor/src/utils/pagination"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCountryService struct {
	FuncCreate func(req *countries.CountryInput) resterror.RestErrorI
	FuncUpdate func(req *countries.UpdateCountryInput, id int64) resterror.RestErrorI
	FuncFind   func(id, asID int64) (*countries.CountryOutput, resterror.RestErrorI)
	FuncList   func(req *countries.ListCountryInput) (*pagination.PaginatedResponse, resterror.RestErrorI)
	FuncDelete func(id int64) resterror.RestErrorI
}

func (m MockCountryService) Create(req *countries.CountryInput) resterror.RestErrorI {
	return m.FuncCreate(req)
}

func (m MockCountryService) Update(req *countries.UpdateCountryInput, id int64) resterror.RestErrorI {
	return m.FuncUpdate(req, id)
}

func (m MockCountryService) Find(id, asID int64) (*countries.CountryOutput, resterror.RestErrorI) {
	return m.FuncFind(id, asID)
}

func (m MockCountryService) List(req *countries.ListCountryInput) (*pagination.PaginatedResponse, resterror.RestErrorI) {
	return m.FuncList(req)
}

func (m MockCountryService) Delete(id int64) resterror.RestErrorI {
	return m.FuncDelete(id)
}

func TestCountryController_Create(t *testing.T) {
	testCases := []struct {
		title          string
		reqBody        io.Reader
		serviceMock    services.CountryServiceI
		expectedStatus int
		expectedRes    string
	}{
		{
			title:          "error validation",
			reqBody:        strings.NewReader(`{}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":{"as_id":["The as id field is required."],"name":["The name field is required."]},"code":400}`,
		},
		{
			title:          "error required name",
			reqBody:        strings.NewReader(`{"as_id":1}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":{"name":["The name field is required."]},"code":400}`,
		},
		{
			title:          "error required as_id",
			reqBody:        strings.NewReader(`{"name":"Englang"}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":{"as_id":["The as id field is required."]},"code":400}`,
		},
		{
			title:          "error invalid active",
			reqBody:        strings.NewReader(`{"name":"England","active":2}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":"Invalid request body.","code":400}`,
		},
		{
			title:   "error CountryService.Create",
			reqBody: strings.NewReader(`{"as_id":1,"name":"England","active":true}`),
			serviceMock: &MockCountryService{
				FuncCreate: func(req *countries.CountryInput) resterror.RestErrorI {
					return resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes:    `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title:   "success",
			reqBody: strings.NewReader(`{"as_id":1,"name":"England","active":true}`),
			serviceMock: &MockCountryService{
				FuncCreate: func(req *countries.CountryInput) resterror.RestErrorI {
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedRes:    `{"message":"SUCCESS","code":201}`,
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

func TestCountryController_Update(t *testing.T) {
	testCases := []struct {
		title          string
		id             string
		reqBody        io.Reader
		serviceMock    services.CountryServiceI
		expectedStatus int
		expectedRes    string
	}{
		{
			title:          "error invalid country id",
			id:             "abc",
			reqBody:        strings.NewReader(`{}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":"INVALID_COUNTRY_ID","code":400}`,
		},
		{
			title:          "error required name",
			id:             "1",
			reqBody:        strings.NewReader(`{}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":{"name":["The name field is required."]},"code":400}`,
		},
		{
			title:          "error invalid active",
			id:             "1",
			reqBody:        strings.NewReader(`{"name":"England","active":2}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":"Invalid request body.","code":400}`,
		},
		{
			title:   "error CountryService.Update",
			id:      "1",
			reqBody: strings.NewReader(`{"name":"England","active":true}`),
			serviceMock: &MockCountryService{
				FuncUpdate: func(req *countries.UpdateCountryInput, id int64) resterror.RestErrorI {
					return resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes:    `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title:   "success",
			id:      "1",
			reqBody: strings.NewReader(`{"name":"England","active":true}`),
			serviceMock: &MockCountryService{
				FuncUpdate: func(req *countries.UpdateCountryInput, id int64) resterror.RestErrorI {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedRes:    `{"message":"SUCCESS","code":200}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "https://localhost:8000/v1/countries"+testCase.id, testCase.reqBody)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			c := utils.GetMockedContext(req, res)
			c.Params = []gin.Param{{Key: "id", Value: testCase.id}}

			services.CountryService = testCase.serviceMock
			CountryController.Update(c)

			assert.Equal(t, testCase.expectedStatus, res.Code)
			assert.Equal(t, testCase.expectedRes, res.Body.String())
		})
	}
}
