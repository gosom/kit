package lib_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gosom/kit/lib"
	"github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
	customValidator := lib.RegisterValidator{
		Tag: "isRizzla",
		Fn: func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			return value == "rizzla"
		},
	}
	lib.SetValidator(customValidator)
	t.Run("TestThatValidatorReturnsErrorIfFieldIsNotRizzla", func(t *testing.T) {
		type Test struct {
			Name string `json:"name" validate:"isRizzla"`
		}
		test := Test{Name: "foo"}
		err := lib.Validate(test)
		require.Error(t, err)
	})
	t.Run("TestThatValidatorReturnsNoErrorIfFieldIsRizzla", func(t *testing.T) {
		type Test struct {
			Name string `json:"name" validate:"isRizzla"`
		}
		test := Test{Name: "rizzla"}
		err := lib.Validate(test)
		require.NoError(t, err)
	})
	t.Run("TestThatValidatorReturnsErrorOnWeakPassword", func(t *testing.T) {
		type Test struct {
			Password string `json:"password" validate:"password"`
		}
		test := Test{Password: "foo"}
		err := lib.Validate(test)
		require.Error(t, err)
	})
	t.Run("TestThatValidatorReturnsNoErrorOnStrongPassword", func(t *testing.T) {
		type Test struct {
			Password string `json:"password" validate:"password"`
		}
		test := Test{Password: "c1!edf0112#"}
		err := lib.Validate(test)
		require.NoError(t, err)
	})
}
