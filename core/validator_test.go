package core_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gosom/kit/core"
	"github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
	customValidator := core.RegisterValidator{
		Tag: "isRizzla",
		Fn: func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			return value == "rizzla"
		},
	}
	core.SetValidator(customValidator)
	t.Run("TestThatValidatorReturnsErrorIfFieldIsNotRizzla", func(t *testing.T) {
		type Test struct {
			Name string `json:"name" validate:"isRizzla"`
		}
		test := Test{Name: "foo"}
		err := core.Validate(test)
		require.Error(t, err)
	})
	t.Run("TestThatValidatorReturnsNoErrorIfFieldIsRizzla", func(t *testing.T) {
		type Test struct {
			Name string `json:"name" validate:"isRizzla"`
		}
		test := Test{Name: "rizzla"}
		err := core.Validate(test)
		require.NoError(t, err)
	})
	t.Run("TestThatValidatorReturnsErrorOnWeakPassword", func(t *testing.T) {
		type Test struct {
			Password string `json:"password" validate:"password"`
		}
		test := Test{Password: "foo"}
		err := core.Validate(test)
		require.Error(t, err)
	})
	t.Run("TestThatValidatorReturnsNoErrorOnStrongPassword", func(t *testing.T) {
		type Test struct {
			Password string `json:"password" validate:"password"`
		}
		test := Test{Password: "c1!edf0112#"}
		err := core.Validate(test)
		require.NoError(t, err)
	})
}
