package services

import (
	"database/sql"
	"fmt"
	"github.com/development-raul/footy-predictor/src/domains/countries"
	"github.com/development-raul/footy-predictor/src/utils/pagination"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
)

type CountryServiceI interface {
	Create(req *countries.CountryInput) resterror.RestErrorI
	Update(req *countries.UpdateCountryInput, id int64) resterror.RestErrorI
	Find(id, asID int64) (*countries.CountryOutput, resterror.RestErrorI)
	List(req *countries.ListCountryInput) (*pagination.PaginatedResponse, resterror.RestErrorI)
	Delete(id int64) resterror.RestErrorI
}

type countryService struct{}

var CountryService CountryServiceI = &countryService{}

func (s *countryService) Create(req *countries.CountryInput) resterror.RestErrorI {
	if err := countries.CountryDao.Create(&countries.Country{
		AsID:   req.AsID,
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

func (s *countryService) Find(id, asID int64) (*countries.CountryOutput, resterror.RestErrorI) {
	var res *countries.CountryOutput
	var err error
	if id != 0 {
		res, err = countries.CountryDao.FindByID(id)
	} else {
		res, err = countries.CountryDao.FindByAsID(asID)
	}

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
	fmt.Println("raul", req.Page, req.PerPage, total)
	res := pagination.GeneratePaginatedResponse(results, req.Page, req.PerPage, total)

	return &res, nil
}

func (s *countryService) Delete(id int64) resterror.RestErrorI {
	if err := countries.CountryDao.Delete(id); err != nil {
		return resterror.NewStandardInternalServerError()
	}
	return nil
}
