package countries

type Country struct {
	ID     int64  `db:"id"`
	Code   string `db:"code"`
	Name   string `db:"name"`
	Flag   string `db:"flag"`
	Active bool   `db:"active"`
}

type CountryInput struct {
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
	ID     int64  `json:"-" form:"-" db:"id"`
	Code   string `json:"code" form:"code" db:"code"`
	Name   string `json:"name" form:"name" db:"name" validate:"required"`
	Flag   string `json:"flag" form:"flag" db:"flag"`
	Active bool   `json:"active" form:"active" db:"active"`
}

type CountryOutput struct {
	ID     int64  `json:"id" db:"id"`
	Code   string `json:"code" db:"code"`
	Name   string `json:"name" db:"name"`
	Flag   string `json:"flag" db:"flag"`
	Active bool   `json:"active" db:"active"`
}
