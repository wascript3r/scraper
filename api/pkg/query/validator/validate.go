package validator

import "github.com/go-playground/validator/v10"

type Validate struct {
	govalidate *validator.Validate
}

func New() *Validate {
	return &Validate{validator.New()}
}

func (v *Validate) RawRequest(s interface{}) error {
	return v.govalidate.Struct(s)
}
