package countries

type Country struct {
	ID     int64  `db:"id"`
	AsID   int64  `db:"as_id"`
	Code   string `db:"code"`
	Name   string `db:"name"`
	Flag   string `db:"flag"`
	Active bool   `db:"active"`
}

type CountryInput struct {
	AsID   int64  `json:"as_id" form:"as_id" validate:"required"`
	Code   string `json:"code" form:"code"`
	Name   string `json:"name" form:"name" validate:"required"`
	Flag   string `json:"flag" form:"flag"`
	Active bool   `json:"active" form:"active"`
}

type ListCountryInput struct {
	Code    string `json:"code" form:"code"`
	Name    string `json:"name" form:"name"`
	Active  bool   `json:"active" form:"active"`
	Order   string `json:"order" form:"order" validate:"omitempty,oneof=desc asc"`
	OrderBy string `json:"order_by" form:"order_by,omitempty" validate:"omitempty,oneof=id code name active"`
	Page    int64  `json:"page" form:"page"`
	PerPage int64  `json:"per_page" form:"per_page"`
}

type UpdateCountryInput struct {
	ID     int64
	Code   string `json:"code" form:"code"`
	Name   string `json:"name" form:"name" validate:"required"`
	Flag   string `json:"flag" form:"flag"`
	Active bool   `json:"active" form:"active"`
}

type CountryOutput struct {
	ID     int64  `json:"id" db:"id"`
	AsID   int64  `json:"as_id" db:"as_id"`
	Code   string `json:"code" db:"code"`
	Name   string `json:"name" db:"name"`
	Flag   string `json:"flag" db:"flag"`
	Active bool   `json:"active" db:"active"`
}
