package countries

import (
	"fmt"
	"github.com/development-raul/footy-predictor/src/clients/mysql/footy_db"
	"github.com/development-raul/footy-predictor/src/utils/helpers"
	"github.com/development-raul/footy-predictor/src/utils/pagination"
	"github.com/development-raul/footy-predictor/src/zlog"
	"strings"
)

type CountryDaoI interface {
	Create(country *Country) error
	Update(country *UpdateCountryInput) error
	FindByID(id int64) (*CountryOutput, error)
	FindByAsID(id int64) (*CountryOutput, error)
	List(req *ListCountryInput) ([]CountryOutput, int64, error)
	Delete(id int64) error
}
type countryDao struct{}

var CountryDao CountryDaoI = &countryDao{}

func (d *countryDao) Create(country *Country) error {
	res, err := footy_db.Client.NamedExec(queryCreate, country)
	if err != nil {
		zlog.Logger.Error("CountryDao Create NamedExec", err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		zlog.Logger.Error("CountryDao Create LastInsertId", err)
		return err
	}
	country.ID = id
	return nil
}

func (d *countryDao) Update(country *UpdateCountryInput) error {
	_, err := footy_db.Client.NamedExec(queryUpdate, country)
	if err != nil {
		zlog.Logger.Error("CountryDao Update NamedExec", err)
		return err
	}
	return nil
}

func (d *countryDao) FindByID(id int64) (*CountryOutput, error) {
	var result CountryOutput

	err := footy_db.Client.Get(&result, queryFindByID, id)
	if err != nil {
		zlog.Logger.Error("CountryDao FindByID Get", err)
		return nil, err
	}
	return &result, nil
}

func (d *countryDao) FindByAsID(id int64) (*CountryOutput, error) {
	var result CountryOutput

	err := footy_db.Client.Get(&result, queryFindByAsID, id)
	if err != nil {
		zlog.Logger.Error("CountryDao FindByAsID Get", err)
		return nil, err
	}
	return &result, nil
}

func (d *countryDao) List(req *ListCountryInput) ([]CountryOutput, int64, error) {
	var results []CountryOutput
	fmt.Println("raul", req.Page, req.PerPage)
	// Create where, limit and order by clauses
	where, args := d.generateListWhereClause(req)
	limit := pagination.GeneratePaginationQuery(req.Page, req.PerPage)
	order := pagination.GeneratePaginationSort("name ASC", req.OrderBy, req.Order)
	query := fmt.Sprintf(queryList, where, order, limit)
	fmt.Println(query, args)
	// Get the records
	err := footy_db.Client.Select(&results, query, args...)
	if err != nil {
		zlog.Logger.Error("CountryDao List Select", err)
		return nil, 0, err
	}

	// Get total records so we can use them for pagination
	total, err := pagination.GetTableTotalRowsArgs(fmt.Sprintf(queryListTotal, where), args...)
	if err != nil {
		zlog.Logger.Error("CountryDao List GetTableTotalRowsArgs", err)
		return nil, 0, err
	}

	return results, total, nil
}

func (d *countryDao) generateListWhereClause(req *ListCountryInput) (string, []interface{}) {
	w := helpers.NewWhere()
	w.AppendWhereAtStart()
	w.Where("true") // add this just in case we do not have any param passed

	if strings.TrimSpace(req.Code) != "" {
		w.Where("code = ?", req.Code)
	}

	if strings.TrimSpace(req.Name) != "" {
		w.CustomWhere(" AND (name LIKE ?)", fmt.Sprintf("%%%s%%", req.Name))
	}

	if req.Active {
		w.Where("active = 1")
	}

	return w.String()
}

func (d *countryDao) Delete(id int64) error {
	_, err := footy_db.Client.Exec(queryDelete, id)
	if err != nil {
		zlog.Logger.Error("CountryDao Delete Exec", err)
		return err
	}
	return nil
}
