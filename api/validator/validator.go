/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package validator

import (
	"errors"
	"regexp"
	"time"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 校验器
type IValidator interface {
	// 注册校验规则
	RegisterValidationRule(tag string, fn validator.Func) error
	// 校验一个结构体
	Valid(a interface{}) error
	// 校验一个字段
	ValidField(a interface{}, tag string) error
}

type Validator struct {
	validateTrans ut.Translator
	validate      *validator.Validate
}

func NewValidator() IValidator {
	zh := zhongwen.New()
	vt, _ := ut.New(zh, zh).GetTranslator("zh")

	validate := validator.New()
	_ = zh_translations.RegisterDefaultTranslations(validate, vt)

	_ = validate.RegisterValidation("regex", validateRegex)
	_ = validate.RegisterValidation("time", validateTime)
	_ = validate.RegisterValidation("date", validateDate)
	return &Validator{
		validateTrans: vt,
		validate:      validate,
	}
}

// 正则匹配
func validateRegex(f validator.FieldLevel) bool {
	compile := f.Param()
	text := f.Field().String()
	return regexp.MustCompile(compile).MatchString(text)
}

// 时间匹配
func validateTime(f validator.FieldLevel) bool {
	layout := f.Param()
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	text := f.Field().String()

	_, err := time.ParseInLocation(layout, text, time.Local)
	return err == nil
}

// 日期匹配
func validateDate(f validator.FieldLevel) bool {
	layout := f.Param()
	if layout == "" {
		layout = "2006-01-02"
	}
	text := f.Field().String()

	_, err := time.ParseInLocation(layout, text, time.Local)
	return err == nil
}

func (v *Validator) RegisterValidationRule(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

func (v *Validator) Valid(a interface{}) error {
	err := v.validate.Struct(a)
	return v.translateValidateErr(err)
}

func (v *Validator) ValidField(a interface{}, tag string) error {
	err := v.validate.Var(a, tag)
	return v.translateValidateErr(err)
}

func (v *Validator) translateValidateErr(err error) error {
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		for _, e := range errs {
			return errors.New(e.Translate(v.validateTrans))
		}
	}
	return nil
}
