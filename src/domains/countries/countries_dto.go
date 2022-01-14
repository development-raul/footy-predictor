package countries

type Country struct {
	ID     int64  `db:"id"`
	Code   string `db:"code"`
	Name   string `db:"name"`
	Flag   string `db:"flag"`
	Active bool   `db:"active"`
}

type CreateCountryInput struct {
	Code   string `json:"code" form:"code"`
	Name   string `json:"name" form:"name" validate:"required"`
	Flag   string `json:"flag" form:"flag"`
	Active int8   `json:"active" form:"active" validate:"omitempty,oneof=0 1"`
}
