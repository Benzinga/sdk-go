package prompt

import validator "github.com/go-playground/validator/v10"

var validate = validator.New()

func NewValidator(validationTag string) func(string) error {
	validationFunc := func(s string) error {
		return validate.Var(s, validationTag)
	}

	return validationFunc
}

func existsValidator() func(string) error {
	return NewValidator("required")
}
