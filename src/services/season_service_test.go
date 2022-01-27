package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/development-raul/footy-predictor/src/clients/restclient"
	"github.com/development-raul/footy-predictor/src/domains/seasons"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

type MockSeasonDao struct {
	FuncCreate func(id int64) error
	FuncFind   func(id int64) (*seasons.Season, error)
	FuncList   func(req *seasons.ListSeasonInput) ([]seasons.Season, error)
	FuncDelete func(id int64) error
}

func (m MockSeasonDao) Create(id int64) error {
	return m.FuncCreate(id)
}
func (m MockSeasonDao) Find(id int64) (*seasons.Season, error) {
	return m.FuncFind(id)
}
func (m MockSeasonDao) List(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
	return m.FuncList(req)
}
func (m MockSeasonDao) Delete(id int64) error {
	return m.FuncDelete(id)
}

func TestSeasonService_Create(t *testing.T) {
	testCases := []struct {
		title         string
		seasonDaoMock seasons.SeasonDaoI
		expectedErr   resterror.RestErrorI
	}{
		{
			title: "error SeasonDao.Create",
			seasonDaoMock: &MockSeasonDao{
				FuncCreate: func(id int64) error {
					return errors.New("error Create")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "success",
			seasonDaoMock: &MockSeasonDao{
				FuncCreate: func(id int64) error {
					return nil
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			seasons.SeasonDao = testCase.seasonDaoMock

			err := SeasonService.Create(1)

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestSeasonService_Find(t *testing.T) {
	testCases := []struct {
		title         string
		id            int64
		seasonDaoMock seasons.SeasonDaoI
		expectedRes   *seasons.Season
		expectedErr   resterror.RestErrorI
	}{
		{
			title: "error SeasonDao.Find",
			id:    1,
			seasonDaoMock: &MockSeasonDao{
				FuncFind: func(id int64) (*seasons.Season, error) {
					return nil, errors.New("error Find")
				},
			},
			expectedRes: nil,
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "error SeasonDao.Find no rows",
			id:    1,
			seasonDaoMock: &MockSeasonDao{
				FuncFind: func(id int64) (*seasons.Season, error) {
					return nil, sql.ErrNoRows
				},
			},
			expectedRes: nil,
			expectedErr: nil,
		},
		{
			title: "success",
			id:    1,
			seasonDaoMock: &MockSeasonDao{
				FuncFind: func(id int64) (*seasons.Season, error) {
					return &seasons.Season{ID: 1}, nil
				},
			},
			expectedRes: &seasons.Season{ID: 1},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			seasons.SeasonDao = testCase.seasonDaoMock

			res, err := SeasonService.Find(testCase.id)

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestSeasonService_List(t *testing.T) {
	testCases := []struct {
		title         string
		seasonDaoMock seasons.SeasonDaoI
		expectedRes   []seasons.Season
		expectedErr   resterror.RestErrorI
	}{
		{
			title: "error SeasonDao.List",
			seasonDaoMock: &MockSeasonDao{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
					return nil, errors.New("error List")
				},
			},
			expectedRes: nil,
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "error SeasonDao.List not found",
			seasonDaoMock: &MockSeasonDao{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
					return nil, sql.ErrNoRows
				},
			},
			expectedRes: nil,
			expectedErr: nil,
		},
		{
			title: "success",
			seasonDaoMock: &MockSeasonDao{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
					return []seasons.Season{{ID: 1}}, nil
				},
			},
			expectedRes: []seasons.Season{{ID: 1}},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			seasons.SeasonDao = testCase.seasonDaoMock

			res, err := SeasonService.List(&seasons.ListSeasonInput{ID: 1, Order: "asc"})

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestSeasonService_Delete(t *testing.T) {
	testCases := []struct {
		title         string
		seasonDaoMock seasons.SeasonDaoI
		expectedErr   resterror.RestErrorI
	}{
		{
			title: "error CountryDao.Delete",
			seasonDaoMock: &MockSeasonDao{
				FuncDelete: func(id int64) error {
					return errors.New("error Delete")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "success",
			seasonDaoMock: &MockSeasonDao{
				FuncDelete: func(id int64) error {
					return nil
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			seasons.SeasonDao = testCase.seasonDaoMock
			err := SeasonService.Delete(1)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestSeasonService_Sync(t *testing.T) {
	os.Setenv("AS_BASE_URL", "http://localhost")
	testCases := []struct {
		title          string
		seasonDaoMock  seasons.SeasonDaoI
		restClientResp *http.Response
		expectedErr    resterror.RestErrorI
	}{
		{
			title: "error SeasonDao.List",
			seasonDaoMock: &MockSeasonDao{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
					return nil, errors.New("error List")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "error api_sports_provider.GetSeasons",
			seasonDaoMock: &MockSeasonDao{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
					return []seasons.Season{{ID: 1}}, nil
				},
			},
			restClientResp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(strings.NewReader(``)),
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "error SeasonDao.Create",
			seasonDaoMock: &MockSeasonDao{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
					return []seasons.Season{{ID: 1}}, nil
				},
				FuncCreate: func(id int64) error {
					return errors.New("error Create")
				},
			},
			restClientResp: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`{
					"get": "leagues/seasons",
					"parameters": [],
					"errors": [],
					"results": 1,
					"paging": {
						"current": 1,
						"total": 1
					},
					"response": [
						2008
					]
				}`)),
			},
			expectedErr: nil,
		},
		{
			title: "success",
			seasonDaoMock: &MockSeasonDao{
				FuncList: func(req *seasons.ListSeasonInput) ([]seasons.Season, error) {
					return []seasons.Season{{ID: 1}}, nil
				},
				FuncCreate: func(id int64) error {
					return nil
				},
			},
			restClientResp: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`{
					"get": "leagues/seasons",
					"parameters": [],
					"errors": [],
					"results": 2,
					"paging": {
						"current": 1,
						"total": 1
					},
					"response": [
						2008,
						2009
					]
				}`)),
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			// Initialization
			restclient.StartMockups()
			restclient.FlushMockups()
			restclient.AddMockup(restclient.Mock{
				Url:        fmt.Sprintf("%s/leagues/seasons", os.Getenv("AS_BASE_URL")),
				HttpMethod: http.MethodGet,
				Response:   testCase.restClientResp,
			})
			seasons.SeasonDao = testCase.seasonDaoMock

			// Execution
			err := SeasonService.Sync()

			// Assertions
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
