package util

import (
	_ "log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	var validator = validator.New()
	validator.RegisterCustomTypeFunc(UtilDateValidateValuer, Date{})
	validator.RegisterCustomTypeFunc(UtilDateValidateValuer, Iso8601{})
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		// validator.RegisterValidation("trimspaces", TrimSpacesValidator)
		return name
	})

	return &Validator{
		validator: validator,
	}
}

func (v *Validator) Validate(s interface{}) error {
	err := v.validator.Struct(s)
	return err
}

// func TrimSpacesValidator(fl validator.FieldLevel) bool {
// 	fieldValue := fl.Field()
// 	if fieldValue.Kind() == reflect.String {
// 		if fieldValue.CanSet() {
// 			trimmedValue := strings.TrimSpace(fieldValue.String())
// 			fieldValue.SetString(trimmedValue)
// 		}
// 	}
// 	return true
// }
