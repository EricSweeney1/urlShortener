package validator

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	v *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{v: validator.New()}
}

func (c *CustomValidator) Validate(i interface{}) error {
	return c.v.Struct(i)
}
