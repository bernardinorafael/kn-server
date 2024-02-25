package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bernardinorafael/gozinho/internal/infra/http/errors"
	"github.com/go-playground/validator/v10"
)

func ValidatorPayload(payload interface{}) *errors.HttpError {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return ""
		}
		return name
	})

	err := v.Struct(payload)
	if err != nil {
		var causes []errors.Fields

		for _, err := range err.(validator.ValidationErrors) {
			cause := errors.Fields{}
			field := err.Field()
			value := err.Value()
			tag := err.Tag()

			switch tag {
			case "required":
				cause.Message = fmt.Sprintf("%s is required", field)
				cause.Field = field
				cause.Value = value
			case "uuid":
				cause.Message = fmt.Sprintf("%s is not a valid uuid", field)
				cause.Field = field
				cause.Value = value
			case "boolean":
				cause.Message = fmt.Sprintf("%s is not a valid boolean", field)
				cause.Field = field
				cause.Value = value
			case "min":
				cause.Message = fmt.Sprintf("%s must be greater than %s", field, err.Param())
				cause.Field = field
				cause.Value = value
			case "max":
				cause.Message = fmt.Sprintf("%s must be less than %s", field, err.Param())
				cause.Field = field
				cause.Value = value
			case "email":
				cause.Message = fmt.Sprintf("%s is not a valid email", field)
				cause.Field = field
				cause.Value = value
			default:
				cause.Message = "invalid field"
				cause.Field = field
				cause.Value = value
			}
			causes = append(causes, cause)
		}
		return errors.NewBadRequestValidationError("some fields are invalid!", causes)
	}
	return nil
}
