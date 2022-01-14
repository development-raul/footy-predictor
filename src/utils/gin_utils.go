package utils

import (
	"fmt"
	"github.com/development-raul/footy-predictor/src/utils/constants"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type CustomUpdate interface {
	Customise()
}

// GinShouldPassAll used to run multiple validation functions using context
// Will return false if any of the functions fail
func GinShouldPassAll(c *gin.Context, funcs ...func(*gin.Context) bool) bool {
	for _, f := range funcs {
		if ok := f(c); !ok {
			return false
		}
	}

	return true
}

// GinShouldBind binds according to content type found in request body, e.g. JSON, XML, Multipart Form, etc
func GinShouldBind(dataPointer interface{}) func(*gin.Context) bool {
	return func(c *gin.Context) bool {
		if err := c.ShouldBind(dataPointer); err != nil {
			restErr := resterror.NewBadRequestError(constants.ErrorInvalidRequestBody)
			c.JSON(restErr.Code(), restErr)
			return false
		}
		return true
	}
}

// GinShouldValidate validates request body using validator v10
func GinShouldValidate(data interface{}) func(*gin.Context) bool {
	return func(c *gin.Context) bool {
		if err := ValidateStruct(data); err != nil {
			c.JSON(http.StatusBadRequest, resterror.ValidationError{
				Error: err,
				Code:  http.StatusBadRequest,
			})
			return false
		}
		// Check for customise method
		val, ok := data.(CustomUpdate)
		if ok {
			val.Customise()
		}
		return true
	}
}

// ValidateStruct validates provided struct using validator v10
func ValidateStruct(obj interface{}) map[string][]string {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})

	// Custom Validation for integers
	_ = v.RegisterValidation("integer", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			_, err := strconv.ParseInt(fl.Field().String(), 10, 64)
			if err != nil {
				return false
			}
		}
		return true
	})

	// Custom Validation for date to be in format YYYY-MM-DD
	_ = v.RegisterValidation("YYYY-MM-DD", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			_, err := time.Parse("2006-01-02", fl.Field().String())
			if err != nil {
				return false
			}
		}
		return true
	})

	_ = v.RegisterValidation("timestamp", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("2006-01-02 15:04:05", fl.Field().String())
		if err != nil {
			return false
		}
		return true
	})

	_ = v.RegisterValidation("futureOnly", func(fl validator.FieldLevel) bool {
		var field time.Time
		// Time with format YYYY-MM-DD hh:mm:ss
		field, err := time.Parse("2006-01-02 15:04:05", fl.Field().String())
		if err != nil {
			// Time with format YYYY-MM-DD
			field, err = time.Parse("2006-01-02", fl.Field().String())
			if err != nil {
				return false
			}
		}

		return field.After(time.Now())
	})

	err := v.Struct(obj)

	if err == nil {
		return nil
	}

	errMap := make(map[string][]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := e.Namespace()
		fieldSplit := strings.Split(field, ".")
		fieldSplit = fieldSplit[1:]
		field = strings.Join(fieldSplit, ".")
		fieldNoUnderscore := strings.ReplaceAll(field, "_", " ")

		switch e.ActualTag() {
		case "len":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must have a length of %v", fieldNoUnderscore, e.Param()))
		case "required":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v field is required.", fieldNoUnderscore))
		case "max":
			errMap[field] = append(errMap[field], fmt.Sprintf("the %v must have a length less than %v", fieldNoUnderscore, e.Param()))
		case "min":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must have a length of at least %v", fieldNoUnderscore, e.Param()))
		case "email":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must be a valid email address.", fieldNoUnderscore))
		case "oneof":
			errMap[field] = append(errMap[field], fmt.Sprintf("The field: '%v' must be one of [%v]", field, e.Param()))
		case "YYYY-MM-DD":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must have the format YYYY-MM-DD", fieldNoUnderscore))
		case "timestamp":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must have the format YYYY-MM-DD hh:mm:ss", fieldNoUnderscore))
		case "futureOnly":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must be a date in the future", fieldNoUnderscore))
		case "integer":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must be a valid integer", fieldNoUnderscore))
		case "UniqueUserEmail":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v has already been taken", fieldNoUnderscore))
		case "ValidSagemakerPerformanceLevel":
			errMap[field] = append(errMap[field], fmt.Sprintf("The selected %v is invalid", fieldNoUnderscore))
		case "e164":
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v must have the format +1234567890", fieldNoUnderscore))
		case "required_if":
			requiredIfParams := strings.Split(e.Param(), " ")
			// Get other fields json tag and replace underscores with spaces
			otherField := strings.ReplaceAll(fieldJSONTag(obj, requiredIfParams[0]), "_", " ")
			if len(requiredIfParams) == 2 {
				errMap[field] = append(errMap[field], fmt.Sprintf("The %v field is required if %v is %v", fieldNoUnderscore, otherField, requiredIfParams[1]))
			} else {
				errMap[field] = append(errMap[field], fmt.Sprintf("The %v field is required with %v", fieldNoUnderscore, otherField))
			}
		case "required_with":
			otherField := strings.ReplaceAll(fieldJSONTag(obj, e.Param()), "_", " ")
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v field is required with %v", fieldNoUnderscore, otherField))
		case "required_without":
			otherField := strings.ReplaceAll(fieldJSONTag(obj, e.Param()), "_", " ")
			errMap[field] = append(errMap[field], fmt.Sprintf("The %v field is required if %v is not present", fieldNoUnderscore, otherField))
		default:
			continue
		}
	}

	return errMap
}

func fieldJSONTag(structVal interface{}, fieldName string) string {
	val := reflect.ValueOf(structVal)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t, _ := val.Type().FieldByName(fieldName)

	switch jsonTag := t.Tag.Get("json"); jsonTag {
	case "-":
		return ""
	case "":
		return ""
	default:
		parts := strings.Split(jsonTag, ",")
		name := parts[0]
		if name == "" {
			name = fieldName
		}
		return name
	}
}

func GetMockedContext(request *http.Request, response *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(response)
	c.Request = request
	return c
}