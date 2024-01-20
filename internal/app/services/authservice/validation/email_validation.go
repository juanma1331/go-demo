package authservice

import (
	"go-demo/internal/app/services"
	"go-demo/internal/app/services/authservice"

	"github.com/go-playground/validator/v10"
)

type uniqueEmailValidator struct {
	UserRepository authservice.AuthUserRepository
}

func (v uniqueEmailValidator) uniqueEmailValidation(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	_, err := v.UserRepository.SelectUserByEmail(email)
	return err != nil
}

func RegisterUniqueEmailValidator(v services.Validator, userRepository authservice.AuthUserRepository) error {
	uniqueEmailValidation := uniqueEmailValidator{
		UserRepository: userRepository,
	}

	return v.RegisterValidation(
		"unique_email",
		uniqueEmailValidation.uniqueEmailValidation,
		"{0} is already taken",
	)
}
