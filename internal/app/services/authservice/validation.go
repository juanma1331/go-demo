package authservice

import (
	"github.com/go-playground/validator/v10"
)

type UniqueEmailValidator struct {
	UserRepository AuthUserRepository
}

const ValidatorUniqueEmailKey = "unique_email"

const ValidatorUniqueEmailErrorMsg = "{0} is already taken"

func (v UniqueEmailValidator) UniqueEmailValidation(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	_, err := v.UserRepository.SelectUserByEmail(email)
	return err != nil
}
