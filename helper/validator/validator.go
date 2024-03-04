package validator

import (
	"fmt"

	httperr "github.com/bernardinorafael/kn-server/helper/error"
	"github.com/go-playground/validator/v10"
)

func Validate(d interface{}) []httperr.ValidationField {
	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(d); err != nil {
		var causes []httperr.ValidationField

		for _, e := range err.(validator.ValidationErrors) {
			cause := httperr.ValidationField{}

			switch e.Tag() {
			case "required":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field %s is required", e.StructField(),
				)
			case "email":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field %s must be a valid email address", e.StructField(),
				)
			case "min":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field %s must have a value greater than %s",
					e.StructField(), e.Param(),
				)
			case "len":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field %s must contain exactly %s numbers",
					e.StructField(), e.Param(),
				)
			case "gte":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field %s must be a positive number", e.StructField(),
				)
			case "max":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field %s must have maximum length of %s",
					e.StructField(), e.Param(),
				)
			default:
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field %s is invalid", e.StructField(),
				)
			}
			causes = append(causes, cause)
		}
		return causes
	}
	return nil
}
