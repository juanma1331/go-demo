package validation

import (
	"github.com/juanma1331/go-demo/internal/auth/app/services"
	"github.com/juanma1331/go-demo/internal/shared"

	"github.com/go-playground/validator/v10"
)

type uniqueEmailValidator struct {
	UserRepository services.AuthUserRepository
}

func (v uniqueEmailValidator) uniqueEmailValidation(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	_, err := v.UserRepository.SelectUserByEmail(email)
	return err != nil
}

func RegisterUniqueEmailValidator(v shared.Validator, userRepository services.AuthUserRepository) error {
	uniqueEmailValidation := uniqueEmailValidator{
		UserRepository: userRepository,
	}

	return v.RegisterValidation(
		"unique_email",
		uniqueEmailValidation.uniqueEmailValidation,
		"{0} is already taken",
	)
}
