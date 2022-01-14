package countries

import (
	"github.com/development-raul/footy-predictor/src/clients/mysql/footy_db"
	"github.com/development-raul/footy-predictor/src/zlog"
)

type CountryDaoI interface {
	Create(country *Country) error
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
