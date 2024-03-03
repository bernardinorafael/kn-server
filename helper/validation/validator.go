package validation

type Validator interface {
	ValidateStruct(d interface{}) error
}
