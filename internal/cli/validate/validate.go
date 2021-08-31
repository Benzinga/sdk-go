package validate

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func NewValidator(validationTag string) func(string) error {
	validationFunc := func(s string) error {
		return validate.Var(s, validationTag)
	}

	return validationFunc
}
