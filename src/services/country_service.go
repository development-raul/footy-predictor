package services

import (
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
)

type CountryServiceI interface {
	Create(req *countries.CreateCountryInput) resterror.RestErrorI
}

type countryService struct{}

var CountryService CountryServiceI = &countryService{}

func (s *countryService) Create(req *countries.CreateCountryInput) resterror.RestErrorI {
	if err := countries.CountryDao.Create(&countries.Country{
		Code:   req.Code,
		Name:   req.Name,
		Flag:   req.Flag,
		Active: req.Active == 1,
	}); err != nil {
		return resterror.NewStandardInternalServerError()
	}
	return nil
}
