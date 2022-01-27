package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/development-raul/footy-predictor/src/clients/restclient"
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/utils/pagination"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

type MockCountryDao struct {
	FuncCreate   func(country *countries.Country) error
	FuncUpdate   func(country *countries.UpdateCountryInput) error
	FuncFindByID func(id int64) (*countries.CountryOutput, error)
	FuncList     func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error)
	FuncDelete   func(id int64) error
}

func (m MockCountryDao) Create(country *countries.Country) error {
	return m.FuncCreate(country)
}
func (m MockCountryDao) Update(country *countries.UpdateCountryInput) error {
	return m.FuncUpdate(country)
}
func (m MockCountryDao) FindByID(id int64) (*countries.CountryOutput, error) {
	return m.FuncFindByID(id)
}
func (m MockCountryDao) List(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
	return m.FuncList(req)
}
func (m MockCountryDao) Delete(id int64) error {
	return m.FuncDelete(id)
}

func TestCountryService_Create(t *testing.T) {
	testCases := []struct {
		title          string
		countryDaoMock countries.CountryDaoI
		expectedErr    resterror.RestErrorI
	}{
		{
			title: "error CountryDao.Create",
			countryDaoMock: &MockCountryDao{
				FuncCreate: func(country *countries.Country) error {
					return errors.New("error Create")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "success",
			countryDaoMock: &MockCountryDao{
				FuncCreate: func(country *countries.Country) error {
					return nil
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			countries.CountryDao = testCase.countryDaoMock
			err := CountryService.Create(&countries.CountryInput{
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: true,
			})
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryService_Update(t *testing.T) {
	testCases := []struct {
		title          string
		countryDaoMock countries.CountryDaoI
		expectedErr    resterror.RestErrorI
	}{
		{
			title: "error CountryDao.FindByID",
			countryDaoMock: &MockCountryDao{
				FuncFindByID: func(id int64) (*countries.CountryOutput, error) {
					return nil, errors.New("error FindByID")
				},
			},
			expectedErr: resterror.NewBadRequestError("INVALID_COUNTRY_ID"),
		},
		{
			title: "error CountryDao.Update",
			countryDaoMock: &MockCountryDao{
				FuncFindByID: func(id int64) (*countries.CountryOutput, error) {
					return &countries.CountryOutput{
						ID:     1,
						Code:   "code",
						Name:   "name",
						Flag:   "flag",
						Active: true,
					}, nil
				},
				FuncUpdate: func(country *countries.UpdateCountryInput) error {
					return errors.New("error Update")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "success",
			countryDaoMock: &MockCountryDao{
				FuncFindByID: func(id int64) (*countries.CountryOutput, error) {
					return &countries.CountryOutput{
						ID:     1,
						Code:   "code",
						Name:   "name",
						Flag:   "flag",
						Active: true,
					}, nil
				},
				FuncUpdate: func(country *countries.UpdateCountryInput) error {
					return nil
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			countries.CountryDao = testCase.countryDaoMock
			err := CountryService.Update(&countries.UpdateCountryInput{
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: true,
			}, 1)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryService_Find(t *testing.T) {
	testCases := []struct {
		title          string
		id             int64
		countryDaoMock countries.CountryDaoI
		expectedRes    *countries.CountryOutput
		expectedErr    resterror.RestErrorI
	}{
		{
			title: "error CountryDao.FindByID",
			id:    1,
			countryDaoMock: &MockCountryDao{
				FuncFindByID: func(id int64) (*countries.CountryOutput, error) {
					return nil, errors.New("error FindByID")
				},
			},
			expectedRes: nil,
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "error CountryDao.FindByID no rows",
			id:    1,
			countryDaoMock: &MockCountryDao{
				FuncFindByID: func(id int64) (*countries.CountryOutput, error) {
					return nil, sql.ErrNoRows
				},
			},
			expectedRes: nil,
			expectedErr: nil,
		},
		{
			title: "success",
			id:    1,
			countryDaoMock: &MockCountryDao{
				FuncFindByID: func(id int64) (*countries.CountryOutput, error) {
					return &countries.CountryOutput{
						ID:     1,
						Code:   "code",
						Name:   "name",
						Flag:   "flag",
						Active: true,
					}, nil
				},
			},
			expectedRes: &countries.CountryOutput{
				ID:     1,
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: true,
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			countries.CountryDao = testCase.countryDaoMock

			res, err := CountryService.Find(testCase.id)

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryService_List(t *testing.T) {
	testCases := []struct {
		title          string
		countryDaoMock countries.CountryDaoI
		expectedRes    *pagination.PaginatedResponse
		expectedTotal  int64
		expectedErr    resterror.RestErrorI
	}{
		{
			title: "error CountryDao.List",
			countryDaoMock: &MockCountryDao{
				FuncList: func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
					return nil, 0, errors.New("error List")
				},
			},
			expectedRes:   nil,
			expectedTotal: 0,
			expectedErr:   resterror.NewStandardInternalServerError(),
		},
		{
			title: "error CountryDao.List not found",
			countryDaoMock: &MockCountryDao{
				FuncList: func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
					return nil, 0, sql.ErrNoRows
				},
			},
			expectedRes: &pagination.PaginatedResponse{
				From:        0,
				Data:        []countries.CountryOutput(nil),
				CurrentPage: 1,
				LastPage:    1,
				PerPage:     10,
				To:          0,
				Total:       0,
			},
			expectedTotal: 0,
			expectedErr:   nil,
		},
		{
			title: "success",
			countryDaoMock: &MockCountryDao{
				FuncList: func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
					return []countries.CountryOutput{
						{
							ID:     1,
							Code:   "code",
							Name:   "name",
							Flag:   "flag",
							Active: true,
						},
					}, 1, nil
				},
			},
			expectedRes: &pagination.PaginatedResponse{
				From: 1,
				Data: []countries.CountryOutput{
					{
						ID:     1,
						Code:   "code",
						Name:   "name",
						Flag:   "flag",
						Active: true,
					},
				},
				CurrentPage: 1,
				LastPage:    1,
				PerPage:     10,
				To:          1,
				Total:       1,
			},
			expectedTotal: 1,
			expectedErr:   nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			countries.CountryDao = testCase.countryDaoMock

			res, err := CountryService.List(&countries.ListCountryInput{
				Code:    "code",
				Name:    "name",
				Active:  true,
				Order:   "name",
				OrderBy: "asc",
				Page:    1,
				PerPage: 10,
			})

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryService_Delete(t *testing.T) {
	testCases := []struct {
		title          string
		countryDaoMock countries.CountryDaoI
		expectedErr    resterror.RestErrorI
	}{
		{
			title: "error CountryDao.Delete",
			countryDaoMock: &MockCountryDao{
				FuncDelete: func(id int64) error {
					return errors.New("error Delete")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "success",
			countryDaoMock: &MockCountryDao{
				FuncDelete: func(id int64) error {
					return nil
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			countries.CountryDao = testCase.countryDaoMock
			err := CountryService.Delete(1)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryService_Sync(t *testing.T) {
	os.Setenv("AS_BASE_URL", "http://localhost")
	testCases := []struct {
		title          string
		countryDaoMock countries.CountryDaoI
		restClientResp *http.Response
		expectedErr    resterror.RestErrorI
	}{
		{
			title: "error CountryDao.List",
			countryDaoMock: &MockCountryDao{
				FuncList: func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
					return nil, 0, errors.New("error List")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "error api_sports_provider.GetCountries",
			countryDaoMock: &MockCountryDao{
				FuncList: func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
					return []countries.CountryOutput{
						{
							ID:     1,
							Code:   "code",
							Name:   "name",
							Flag:   "flag",
							Active: true,
						},
					}, 1, nil
				},
			},
			restClientResp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(strings.NewReader(``)),
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "error CountryDao.Create",
			countryDaoMock: &MockCountryDao{
				FuncList: func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
					return []countries.CountryOutput{
						{
							ID:     1,
							Code:   "code",
							Name:   "name",
							Flag:   "flag",
							Active: true,
						},
					}, 1, nil
				},
				FuncCreate: func(country *countries.Country) error {
					return errors.New("error Create")
				},
			},
			restClientResp: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`{
					"get": "countries",
					"parameters": [],
					"errors": [],
					"results": 1,
					"paging": {
						"current": 1,
						"total": 1
					},
					"response": [
						{
							"name": "Albania",
							"code": "AL",
							"flag": "flag_url"
						}
					]
				}`)),
			},
			expectedErr: nil,
		},
		{
			title: "success",
			countryDaoMock: &MockCountryDao{
				FuncList: func(req *countries.ListCountryInput) ([]countries.CountryOutput, int64, error) {
					return []countries.CountryOutput{
						{
							ID:     1,
							Code:   "code",
							Name:   "name",
							Flag:   "flag",
							Active: true,
						},
					}, 1, nil
				},
				FuncCreate: func(country *countries.Country) error {
					return nil
				},
			},
			restClientResp: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`{
					"get": "countries",
					"parameters": [],
					"errors": [],
					"results": 163,
					"paging": {
						"current": 1,
						"total": 1
					},
					"response": [
						{
							"name": "Albania",
							"code": "AL",
							"flag": "flag_url"
						},
						{
							"name": "name",
							"code": "code",
							"flag": "flag"
						}
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
				Url:        fmt.Sprintf("%s/countries", os.Getenv("AS_BASE_URL")),
				HttpMethod: http.MethodGet,
				Response:   testCase.restClientResp,
			})
			countries.CountryDao = testCase.countryDaoMock

			// Execution
			err := CountryService.Sync()

			// Assertions
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
