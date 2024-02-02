package shared

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	Struct(interface{}) (*ValidationErrors, error)
	RegisterValidation(string, validator.Func, string) error
}

type ValidationErrors map[string][]string

type playgroundValidator struct {
	v *validator.Validate
	t ut.Translator
}

func NewPlaygroundValidator(v *validator.Validate) playgroundValidator {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	t, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v, t)

	return playgroundValidator{
		v: v,
		t: t,
	}
}

func (pv playgroundValidator) Struct(i interface{}) (*ValidationErrors, error) {
	err := pv.v.Struct(i)
	return pv.handleValidationError(err)
}

func (pv *playgroundValidator) RegisterValidation(tag string, fn validator.Func, translation string) error {
	if err := pv.registerValidatorFunc(tag, fn); err != nil {
		return err
	}

	return pv.registerValidatorTranslation(tag, translation)
}

func (pv playgroundValidator) handleValidationError(err error) (*ValidationErrors, error) {
	if err == nil {
		return nil, nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return nil, err
	}

	if valErrs, ok := err.(validator.ValidationErrors); ok {

		fmt.Printf("%+v\n", valErrs)
		return pv.convertValidationErrors(valErrs), nil
	}

	return nil, err
}

func (pv playgroundValidator) convertValidationErrors(valErrs validator.ValidationErrors) *ValidationErrors {
	validationErrors := make(ValidationErrors)

	for _, err := range valErrs {
		field := err.Field()
		message := err.Translate(pv.t)

		validationErrors[field] = append(validationErrors[field], message)
	}

	return &validationErrors
}

func (pv *playgroundValidator) registerValidatorFunc(tag string, fn validator.Func) error {
	return pv.v.RegisterValidation(tag, fn)
}

func (pv *playgroundValidator) registerValidatorTranslation(tag string, translation string) error {
	return pv.v.RegisterTranslation(tag, pv.t, func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}, pv.translateFunc)
}

func (pv *playgroundValidator) translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(fe.Tag(), fe.Field())
	return t
}
