package validator

import "github.com/go-playground/validator/v10"

type Validate struct {
	dateTimeFormat string
	govalidate     *validator.Validate
}

func New(dateTimeFormat string) *Validate {
	goV := validator.New()

	r := newRules(dateTimeFormat)
	r.attachTo(goV)

	return &Validate{dateTimeFormat, goV}
}

func (v *Validate) RawRequest(s interface{}) error {
	return v.govalidate.Struct(s)
}

func (v *Validate) GetDateTimeFormat() string {
	return v.dateTimeFormat
}
