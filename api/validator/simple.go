/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/21
   Description :
-------------------------------------------------
*/

package validator

import (
	"github.com/go-playground/validator/v10"
)

var defaultValidator IValidator

func init() {
	defaultValidator = NewValidator()
}

func RegisterValidationRule(tag string, fn validator.Func) error {
	return defaultValidator.RegisterValidationRule(tag, fn)
}

func Valid(a interface{}) error {
	return defaultValidator.Valid(a)
}

func ValidField(a interface{}, tag string) error {
	return defaultValidator.ValidField(a, tag)
}
