package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(i interface{}) error {
	return validate.Struct(i)
}
