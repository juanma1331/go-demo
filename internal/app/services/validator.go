package services

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(interface{}) (*ValidationErrors, error)
}

type ValidationErrors map[string][]string

type playgroundValidator struct {
	v *validator.Validate
}

func NewPlaygroundValidator() playgroundValidator {
	return playgroundValidator{
		v: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (pv playgroundValidator) Struct(i interface{}) (*ValidationErrors, error) {
	err := pv.v.Struct(i)
	return pv.handleValidationError(err)
}

func (pv playgroundValidator) handleValidationError(err error) (*ValidationErrors, error) {
	if err == nil {
		return nil, nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return nil, err
	}

	if valErrs, ok := err.(validator.ValidationErrors); ok {

		return pv.convertValidationErrors(valErrs), nil
	}

	return nil, err
}

func (pv playgroundValidator) convertValidationErrors(valErrs validator.ValidationErrors) *ValidationErrors {
	validationErrors := make(ValidationErrors)

	for _, err := range valErrs {
		field := err.Field()
		message := err.Tag()

		validationErrors[field] = append(validationErrors[field], message)
	}

	return &validationErrors
}
