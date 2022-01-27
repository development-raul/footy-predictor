package services

import (
	"database/sql"
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/providers/api_sports_provider"
	"github.com/development-raul/footy-predictor/src/utils/pagination"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/development-raul/footy-predictor/src/zlog"
)

type CountryServiceI interface {
	Create(req *countries.CountryInput) resterror.RestErrorI
	Update(req *countries.UpdateCountryInput, id int64) resterror.RestErrorI
	Find(id int64) (*countries.CountryOutput, resterror.RestErrorI)
	List(req *countries.ListCountryInput) (*pagination.PaginatedResponse, resterror.RestErrorI)
	Delete(id int64) resterror.RestErrorI
	Sync() resterror.RestErrorI
}

type countryService struct{}

var CountryService CountryServiceI = &countryService{}

func (s *countryService) Create(req *countries.CountryInput) resterror.RestErrorI {
	if err := countries.CountryDao.Create(&countries.Country{
		Code:   req.Code,
		Name:   req.Name,
		Flag:   req.Flag,
		Active: req.Active,
	}); err != nil {
		return resterror.NewStandardInternalServerError()
	}
	return nil
}

func (s *countryService) Update(req *countries.UpdateCountryInput, id int64) resterror.RestErrorI {
	// Check if the country already exists
	country, err := countries.CountryDao.FindByID(id)
	if err != nil {
		return resterror.NewBadRequestError("INVALID_COUNTRY_ID")
	}

	// Set the ID and update the records
	req.ID = country.ID
	if err := countries.CountryDao.Update(req); err != nil {
		return resterror.NewStandardInternalServerError()
	}
	return nil
}

func (s *countryService) Find(id int64) (*countries.CountryOutput, resterror.RestErrorI) {
	res, err := countries.CountryDao.FindByID(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, resterror.NewStandardInternalServerError()
	}
	return res, nil
}

func (s *countryService) List(req *countries.ListCountryInput) (*pagination.PaginatedResponse, resterror.RestErrorI) {
	results, total, err := countries.CountryDao.List(req)
	if err != nil && err != sql.ErrNoRows {
		return nil, resterror.NewStandardInternalServerError()
	}

	res := pagination.GeneratePaginatedResponse(results, req.Page, req.PerPage, total)

	return &res, nil
}

func (s *countryService) Delete(id int64) resterror.RestErrorI {
	if err := countries.CountryDao.Delete(id); err != nil {
		return resterror.NewStandardInternalServerError()
	}
	return nil
}

func (s *countryService) Sync() resterror.RestErrorI {
	zlog.Logger.Info("Sync Countries Start")
	// Get existing countries - set a high pagination, so we can be sure we are getting all in one go
	filters := countries.ListCountryInput{PerPage: 999}
	results, total, err := countries.CountryDao.List(&filters)
	if err != nil && err != sql.ErrNoRows {
		return resterror.NewStandardInternalServerError()
	}
	// Create a map with existing counties, so we can easily identify the already existing countries
	existingCountries := make(map[string]string, total)
	for _, v := range results {
		existingCountries[v.Name] = v.Code
	}
	// Get the list of countries from API Sports
	res, apiErr := api_sports_provider.GetCountries()
	if apiErr != nil {
		return resterror.NewStandardInternalServerError()
	}

	for _, country := range res {
		// Check if the country already exists
		if _, exists := existingCountries[country.Name]; exists {
			continue
		}
		// Create the country if it does not exist
		err := countries.CountryDao.Create(&countries.Country{
			Code:   country.Code,
			Name:   country.Name,
			Flag:   country.Flag,
			Active: true,
		})
		if err != nil {
			zlog.Logger.Warn("could not create country: ", country.Name)
			continue
		}
		zlog.Logger.Info("created new country: ", country.Name)
	}
	zlog.Logger.Info("Sync Countries End")
	return nil
}
