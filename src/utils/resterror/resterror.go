package resterror

import (
"encoding/json"
"errors"
	"github.com/development-raul/footy-predictor/src/utils/constants"
	"net/http"
)

type RestErrorI interface {
	Error() interface{}
	Code() int
}

type restError struct {
	Err        interface{} `json:"error"`
	StatusCode int         `json:"code"`
}

type ValidationError struct {
	Error interface{} `json:"error"`
	Code  int         `json:"code"`
}

func (i *restError) Error() interface{} {
	return i.Err
}

func (i *restError) Code() int {
	return i.StatusCode
}

func NewConflictError(msg string) RestErrorI {
	return &restError{
		Err:        msg,
		StatusCode: http.StatusConflict,
	}
}

func NewUnauthorizedError(msg string) RestErrorI {
	return &restError{
		Err:        msg,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewBadRequestError(msg string) RestErrorI {
	return &restError{
		Err:        msg,
		StatusCode: http.StatusBadRequest,
	}
}

func NewUnprocessableEntityError(msg string) RestErrorI {
	return &restError{
		Err:        msg,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

func NewNotFoundError(msg string) RestErrorI {
	return &restError{
		Err:        msg,
		StatusCode: http.StatusNotFound,
	}
}

func NewInternalServerError(msg string) RestErrorI {
	return &restError{
		Err:        msg,
		StatusCode: http.StatusInternalServerError,
	}
}

func NewStandardInternalServerError() RestErrorI {
	return &restError{
		Err:        constants.ServerError,
		StatusCode: http.StatusInternalServerError,
	}
}

func NewCustomError(e interface{}, code int) RestErrorI {
	return &restError{
		Err:        e,
		StatusCode: code,
	}
}

func NewForbiddenError(msg string) RestErrorI {
	return &restError{
		Err:        msg,
		StatusCode: http.StatusForbidden,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErrorI, error) {
	var apiErr restError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid error json response")
	}
	return &apiErr, nil
}
