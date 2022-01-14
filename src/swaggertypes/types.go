package swaggertypes

type StandardValidationError struct {
	Error struct {
		FieldName []string `json:"field_name" example:"example error"`
	} `json:"error"`
	Code int `json:"code" example:"400"`
}
type StandardInternalServerError struct {
	Error string `json:"error" example:"Server Error"`
	Code  int    `json:"code" example:"500"`
}
type StandardBadRequestError struct {
	Error string `json:"error" example:"Bad Request"`
	Code  int    `json:"code" example:"400"`
}
type StandardNotFoundError struct {
	Error string `json:"error" example:"Not found"`
	Code  int    `json:"code" example:"404"`
}

type NoErrorString struct {
	Message string `json:"message"`
	Code    int    `json:"code" example:"200"`
}

type NoErrorData struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

type NoErrorI struct {
	Code int `json:"code" example:"200"`
}
type StandardUnauthorisedError struct {
	Error string `json:"error" example:"INVALID_USER_AUTHENTICATION"`
	Code  int    `json:"code" example:"401"`
}

type PaginatedData struct {
	Data struct {
		From        int64       `json:"from"`
		Data        interface{} `json:"data"`
		CurrentPage int64       `json:"current_page"`
		LastPage    int64       `json:"last_page"`
		PerPage     int64       `json:"per_page,string"`
		To          int64       `json:"to"`
		Total       int64       `json:"total"`
	} `json:"data"`
	Code int `json:"code"`
}
