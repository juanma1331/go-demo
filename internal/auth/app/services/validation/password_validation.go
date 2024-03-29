package validation

import (
	"unicode"

	"github.com/juanma1331/go-demo/internal/shared"

	"github.com/go-playground/validator/v10"
)

func RegisterPasswordValidations(v shared.Validator) error {
	if err := v.RegisterValidation("lowercase", func(fl validator.FieldLevel) bool {
		return containsRune(fl.Field().String(), isLetterLowercase)
	}, "At least one lowercase letter required"); err != nil {
		return err
	}

	if err := v.RegisterValidation("uppercase", func(fl validator.FieldLevel) bool {
		return containsRune(fl.Field().String(), isLetterUppercase)
	}, "At least one uppercase letter required"); err != nil {
		return err
	}

	if err := v.RegisterValidation("number", func(fl validator.FieldLevel) bool {
		return containsRune(fl.Field().String(), isDigit)
	}, "At least one number required"); err != nil {
		return err
	}

	if err := v.RegisterValidation("special", func(fl validator.FieldLevel) bool {
		return containsRune(fl.Field().String(), isSpecialCharacter)
	}, "At least one special character required"); err != nil {
		return err
	}

	return nil
}

func isLetterLowercase(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}

func isLetterUppercase(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsUpper(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isSpecialCharacter(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r)
}

func containsRune(s string, testFunc func(rune) bool) bool {
	for _, r := range s {
		if testFunc(r) {
			return true
		}
	}
	return false
}
