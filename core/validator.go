package core

import (
	"sync"

	"github.com/go-playground/validator/v10"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var validate *validator.Validate
var validationOnce sync.Once

type RegisterValidator struct {
	Tag string
	Fn  validator.Func
}

func SetValidator(validators ...RegisterValidator) {
	validationOnce.Do(func() {
		passwordValidator := RegisterValidator{
			Tag: "password",
			Fn:  validatePassword,
		}
		validatorsToAdd := []RegisterValidator{passwordValidator}
		for _, v := range validators {
			if v.Tag != "" && v.Fn != nil && v.Tag != "password" {
				validatorsToAdd = append(validatorsToAdd, v)
			}
		}
		validate = validator.New()
		for _, v := range validatorsToAdd {
			validate.RegisterValidation(v.Tag, v.Fn)
		}
	})
}

func NewValidator(validators ...RegisterValidator) *validator.Validate {
	return validate
}

func Validate(v interface{}, skipFields ...string) error {
	SetValidator()
	return validate.StructExcept(v, skipFields...)
}

func validatePassword(fl validator.FieldLevel) bool {
	const minEntropyBits = 50
	err := passwordvalidator.Validate(fl.Field().String(), minEntropyBits)
	return err == nil
}
