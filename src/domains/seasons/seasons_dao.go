package seasons

import (
	"fmt"
	"github.com/development-raul/footy-predictor/src/clients/mysql/footy_db"
	"github.com/development-raul/footy-predictor/src/utils/helpers"
	"github.com/development-raul/footy-predictor/src/utils/pagination"
	"github.com/development-raul/footy-predictor/src/zlog"
)

type SeasonDaoI interface {
	Create(id int64) error
	Find(id int64) (*Season, error)
	List(req *ListSeasonInput) ([]Season, error)
	Delete(id int64) error
}

type seasonDao struct{}

var SeasonDao SeasonDaoI = &seasonDao{}

func (d *seasonDao) Create(id int64) error {
	_, err := footy_db.Client.Exec(queryCreate, id)
	if err != nil {
		zlog.Logger.Error("SeasonDao Create Exec", err)
		return err
	}
	return nil
}

func (d *seasonDao) Find(id int64) (*Season, error) {
	var result Season

	err := footy_db.Client.Get(&result, queryFind, id)
	if err != nil {
		zlog.Logger.Error("SeasonDao Find Get", err)
		return nil, err
	}
	return &result, nil
}

func (d *seasonDao) List(req *ListSeasonInput) ([]Season, error) {
	var results []Season

	// Create where, and order by clauses
	where, args := d.generateListWhereClause(req)
	order := pagination.GeneratePaginationSort("id ASC", "id", req.Order)
	query := fmt.Sprintf(queryList, where, order)

	// Get the records
	err := footy_db.Client.Select(&results, query, args...)
	if err != nil {
		zlog.Logger.Error("SeasonDao List Select", err)
		return nil, err
	}

	return results, nil
}

func (d *seasonDao) generateListWhereClause(req *ListSeasonInput) (string, []interface{}) {
	w := helpers.NewWhere()
	w.AppendWhereAtStart()
	w.Where("true") // add this just in case we do not have any param passed

	if req.ID != 0 {
		w.Where("id = ?", req.ID)
	}
	return w.String()
}

func (d *seasonDao) Delete(id int64) error {
	_, err := footy_db.Client.Exec(queryDelete, id)
	if err != nil {
		zlog.Logger.Error("SeasonDao Delete Exec", err)
		return err
	}
	return nil
}
