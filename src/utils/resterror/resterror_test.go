package resterror

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRestError_Error(t *testing.T) {
	jsonTest(t, NewConflictError("test msg"), "test msg", 409)

	jsonTest(t, NewUnauthorizedError("test msg"), "test msg", 401)

	jsonTest(t, NewBadRequestError("test msg"), "test msg", 400)

	jsonTest(t, NewUnprocessableEntityError("test msg"), "test msg", 422)

	jsonTest(t, NewNotFoundError("test msg"), "test msg", 404)

	jsonTest(t, NewInternalServerError("test msg"), "test msg", 500)

	jsonTest(t, NewStandardInternalServerError(), "Something went wrong. Please try again later.", 500)

	jsonTest(t, NewCustomError("test msg", 202), "test msg", 202)

	jsonTest(t, NewForbiddenError("test msg"), "test msg", 403)
}

func jsonTest(t *testing.T, e RestErrorI, msg string, code int) {
	b, err := json.Marshal(e)
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf(`{"error":"%v","code":%v}`, msg, code), string(b))
	assert.Equal(t, e.Error(), msg)
	assert.Equal(t, e.Code(), code)
}

func TestNewRestErrorFromBytes(t *testing.T) {
	r1, e1 := NewRestErrorFromBytes([]byte(`{"error":"test error","code":500}`))
	assert.Nil(t, e1)
	assert.Equal(t, r1.Error(), "test error")
	assert.Equal(t, r1.Code(), 500)

	r2, e2 := NewRestErrorFromBytes([]byte(`{"error":true,"code":"abc"}`))
	if e2 != nil {
		assert.Equal(t, "invalid error json response", e2.Error())
	}
	assert.Nil(t, r2)
}
