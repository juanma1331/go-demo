package services

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(interface{}) ([]ValidationError, error)
	Var(field interface{}, tag string) ([]ValidationError, error)
}

type playgroundValidator struct {
	v *validator.Validate
}

type ValidationError struct {
	Field   string
	Message string
}

func NewPlaygroundValidator() playgroundValidator {
	return playgroundValidator{
		v: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (pv playgroundValidator) Struct(i interface{}) ([]ValidationError, error) {
	err := pv.v.Struct(i)
	return pv.handleValidationError(err)
}

func (pv playgroundValidator) Var(field interface{}, tag string) ([]ValidationError, error) {
	err := pv.v.Var(field, tag)
	return pv.handleValidationError(err)
}

func (pv playgroundValidator) handleValidationError(err error) ([]ValidationError, error) {
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

func (pv playgroundValidator) convertValidationErrors(valErrs validator.ValidationErrors) []ValidationError {
	var validationErrors []ValidationError
	for _, v := range valErrs {
		ve := ValidationError{
			Field:   v.Field(),
			Message: v.Error(),
		}
		validationErrors = append(validationErrors, ve)
	}
	return validationErrors
}
