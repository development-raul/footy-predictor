package seasons

type Season struct {
	ID int64 `json:"id" form:"id" db:"id" validate:"required"`
}

type ListSeasonInput struct {
	ID    int64  `json:"id" form:"id"`
	Order string `json:"order" form:"order" validate:"omitempty,oneof=desc asc"`
}
