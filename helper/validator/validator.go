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
				cause.Msg = "is a required field"
			case "email":
				cause.Field = e.StructField()
				cause.Msg = "must be a valid email address"
			case "min":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field `%s` must have the length greater than %s",
					e.StructField(), e.Param(),
				)
			case "max":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field `%s` must have a maximum of  %s characters",
					e.StructField(), e.Param(),
				)
			case "len":
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf(
					"the field `%s` must have a length %s characters",
					e.StructField(), e.Param(),
				)
			case "numeric":
				cause.Field = e.StructField()
				cause.Msg = "only numeric digits are allowed"
			default:
				cause.Field = e.StructField()
				cause.Msg = fmt.Sprintf("the field `%s` is invalid", e.StructField())
			}
			causes = append(causes, cause)
		}
		return causes
	}
	return nil
}
