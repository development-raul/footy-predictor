package controllers

import (
	"github.com/development-raul/footy-predictor/src/domains/seasons"
	"github.com/development-raul/footy-predictor/src/services"
	"github.com/development-raul/footy-predictor/src/utils"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockSeasonService struct {
	FuncCreate func(id int64) resterror.RestErrorI
	FuncFind   func(id int64) (*seasons.Season, resterror.RestErrorI)
	FuncList   func(req *seasons.ListSeasonInput) ([]seasons.Season, resterror.RestErrorI)
	FuncDelete func(id int64) resterror.RestErrorI
	FuncSync   func() resterror.RestErrorI
}

func (m MockSeasonService) Create(id int64) resterror.RestErrorI {
	return m.FuncCreate(id)
}
func (m MockSeasonService) Find(id int64) (*seasons.Season, resterror.RestErrorI) {
	return m.FuncFind(id)
}
func (m MockSeasonService) List(req *seasons.ListSeasonInput) ([]seasons.Season, resterror.RestErrorI) {
	return m.FuncList(req)
}
func (m MockSeasonService) Delete(id int64) resterror.RestErrorI {
	return m.FuncDelete(id)
}
func (m MockSeasonService) Sync() resterror.RestErrorI {
	return m.FuncSync()
}

func TestSeasonController_Create(t *testing.T) {
	testCases := []struct {
		title          string
		reqBody        io.Reader
		serviceMock    services.SeasonServiceI
		expectedStatus int
		expectedRes    string
	}{
		{
			title:          "error required id",
			reqBody:        strings.NewReader(`{}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":{"id":["The id field is required."]},"code":400}`,
		},
		{
			title:          "error invalid id",
			reqBody:        strings.NewReader(`{"id":"test"}`),
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":"Invalid request body.","code":400}`,
		},
		{
			title:   "error SeasonService.Create",
			reqBody: strings.NewReader(`{"id":1}`),
			serviceMock: &MockSeasonService{
				FuncCreate: func(id int64) resterror.RestErrorI {
					return resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes:    `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title:   "success",
			reqBody: strings.NewReader(`{"id":1}`),
			serviceMock: &MockSeasonService{
				FuncCreate: func(id int64) resterror.RestErrorI {
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedRes:    `{"message":"SUCCESS","code":201}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "https://localhost:8000/v1/seasons", testCase.reqBody)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			c := utils.GetMockedContext(req, res)

			services.SeasonService = testCase.serviceMock
			SeasonController.Create(c)

			assert.Equal(t, testCase.expectedStatus, res.Code)
			assert.Equal(t, testCase.expectedRes, res.Body.String())
		})
	}
}

func TestSeasonController_Find(t *testing.T) {
	testCases := []struct {
		title          string
		id             string
		serviceMock    services.SeasonServiceI
		expectedStatus int
		expectedRes    string
	}{
		{
			title:          "error invalid season id",
			id:             "abc",
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":"INVALID_SEASON_ID","code":400}`,
		},
		{
			title: "error SeasonService.Find",
			id:    "1",
			serviceMock: &MockSeasonService{
				FuncFind: func(id int64) (*seasons.Season, resterror.RestErrorI) {
					return nil, resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes:    `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title: "success",
			id:    "1",
			serviceMock: &MockSeasonService{
				FuncFind: func(id int64) (*seasons.Season, resterror.RestErrorI) {
					return &seasons.Season{ID: 1}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedRes:    `{"data":{"id":1},"code":200}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "https://localhost:8000/v1/seasons"+testCase.id, nil)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			c := utils.GetMockedContext(req, res)
			c.Params = []gin.Param{{Key: "id", Value: testCase.id}}

			services.SeasonService = testCase.serviceMock
			SeasonController.Find(c)

			assert.Equal(t, testCase.expectedStatus, res.Code)
			assert.Equal(t, testCase.expectedRes, res.Body.String())
		})
	}
}

func TestSeasonController_List(t *testing.T) {
	testCases := []struct {
		title          string
		query          string
		serviceMock    services.SeasonServiceI
		expectedStatus int
		expectedRes    string
	}{
		{
			title:          "error validation invalid order",
			query:          "?order=test",
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":{"order":["The field: 'order' must be one of [desc asc]"]},"code":400}`,
		},
		{
			title:          "error validation invalid id",
			query:          "?id=test",
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":"Invalid request body.","code":400}`,
		},
		{
			title: "error SeasonService.List",
			query: "?order=asc&order_by=id",
			serviceMock: &MockSeasonService{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, resterror.RestErrorI) {
					return nil, resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes:    `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title: "success",
			query: "?order=asc&order_by=id",
			serviceMock: &MockSeasonService{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, resterror.RestErrorI) {
					return []seasons.Season{{ID: 1}}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedRes:    `{"data":[{"id":1}],"code":200}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "https://localhost:8000/v1/seasons"+testCase.query, nil)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			c := utils.GetMockedContext(req, res)

			services.SeasonService = testCase.serviceMock
			SeasonController.List(c)

			assert.Equal(t, testCase.expectedStatus, res.Code)
			assert.Equal(t, testCase.expectedRes, res.Body.String())
		})
	}
}

func TestSeasonController_Delete(t *testing.T) {
	testCases := []struct {
		title          string
		id             string
		serviceMock    services.SeasonServiceI
		expectedStatus int
		expectedRes    string
	}{
		{
			title:          "error invalid season id",
			id:             "abc",
			serviceMock:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedRes:    `{"error":"INVALID_SEASON_ID","code":400}`,
		},
		{
			title: "error SeasonService.Delete",
			id:    "1",
			serviceMock: &MockSeasonService{
				FuncDelete: func(id int64) resterror.RestErrorI {
					return resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes:    `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title: "success",
			id:    "1",
			serviceMock: &MockSeasonService{
				FuncDelete: func(id int64) resterror.RestErrorI {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedRes:    `{"message":"SUCCESS","code":200}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "https://localhost:8000/v1/seasons"+testCase.id, nil)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			c := utils.GetMockedContext(req, res)
			c.Params = []gin.Param{{Key: "id", Value: testCase.id}}

			services.SeasonService = testCase.serviceMock
			SeasonController.Delete(c)

			assert.Equal(t, testCase.expectedStatus, res.Code)
			assert.Equal(t, testCase.expectedRes, res.Body.String())
		})
	}
}

func TestSeasonController_Sync(t *testing.T) {
	testCases := []struct {
		title          string
		serviceMock    services.SeasonServiceI
		expectedStatus int
		expectedRes    string
	}{
		{
			title: "error SeasonService.Sync",
			serviceMock: &MockSeasonService{
				FuncSync: func() resterror.RestErrorI {
					return resterror.NewStandardInternalServerError()
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedRes:    `{"error":"Something went wrong. Please try again later.","code":500}`,
		},
		{
			title: "success",
			serviceMock: &MockSeasonService{
				FuncSync: func() resterror.RestErrorI {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedRes:    `{"message":"SUCCESS","code":200}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "https://localhost:8000/v1/seasons/sync", nil)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			c := utils.GetMockedContext(req, res)

			services.SeasonService = testCase.serviceMock
			SeasonController.Sync(c)

			assert.Equal(t, testCase.expectedStatus, res.Code)
			assert.Equal(t, testCase.expectedRes, res.Body.String())
		})
	}
}
