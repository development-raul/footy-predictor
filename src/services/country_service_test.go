package services

import (
	"errors"
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockCountryDao struct {
	FuncCreate func(country *countries.Country) error
}

func (m MockCountryDao) Create(country *countries.Country) error {
	return m.FuncCreate(country)
}

func TestCountryService_Create(t *testing.T) {
	testCases := []struct{
		title string
		countryDaoMock countries.CountryDaoI
		expectedErr resterror.RestErrorI
	}{
		{
			title: "error CountryDao.Create",
			countryDaoMock: &MockCountryDao{
				func(country *countries.Country) error {
					return errors.New("test create")
				},
			},
			expectedErr: resterror.NewStandardInternalServerError(),
		},
		{
			title: "success",
			countryDaoMock: &MockCountryDao{
				func(country *countries.Country) error {
					return nil
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			countries.CountryDao = testCase.countryDaoMock
			err := CountryService.Create(&countries.CreateCountryInput{
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: 1,
			})
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}