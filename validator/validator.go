package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErrors = validator.ValidationErrors

type FieldViolation struct {
	Field     string `json:"field"`
	Tag       string `json:"tag"`
	Condition string `json:"condition"`
}

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("mapstructure"), ",", 2)[0]

		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		}

		if name == "" || name == "-" {
			return fld.Name
		}

		return name
	})

	return &Validator{v: v}
}

func (val *Validator) Validate(s any) error {
	return val.v.Struct(s)
}

func (val *Validator) ValidateWithErrors(s any) ([]FieldViolation, error) {
	err := val.v.Struct(s)
	if err == nil {
		return nil, nil
	}

	ve, ok := err.(ValidationErrors)
	if !ok {
		return nil, err
	}

	violations := make([]FieldViolation, len(ve))
	for i, fe := range ve {
		violations[i] = FieldViolation{
			Field:     fe.Field(),
			Tag:       fe.Tag(),
			Condition: fe.Param(),
		}
	}
	return violations, err
}

func (val *Validator) RegisterCustom(tag string, fn validator.Func, callIfNull ...bool) error {
	return val.v.RegisterValidation(tag, fn, callIfNull...)
}

func (val *Validator) Engine() *validator.Validate {
	return val.v
}
