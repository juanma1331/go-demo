package authservice

import (
	"fmt"
	"go-demo/internal/app/services"
	"go-demo/internal/domain"
	"net/http"
)

type AuthServiceParams struct {
	UserRepository  AuthUserRepository
	PasswordManager PasswordManager
	SessionStore    SessionStore
	Validator       services.Validator
}

type authService struct {
	userRepo        AuthUserRepository
	passwordManager PasswordManager
	sessionStore    SessionStore
	validator       services.Validator
}

func NewAuthService(params AuthServiceParams) *authService {
	return &authService{
		userRepo:        params.UserRepository,
		passwordManager: params.PasswordManager,
		sessionStore:    params.SessionStore,
		validator:       params.Validator,
	}
}

func (as *authService) Register(i RegisterInput) (RegisterOutput, error) {
	output := RegisterOutput{}

	valErrs, err := as.validator.Struct(i)
	if err != nil {
		return output, err
	}

	if valErrs != nil {
		output.ValidationErrors = valErrs
		return output, nil
	}

	hashedPassword, err := as.passwordManager.GenerateFromPassword(i.Password)
	if err != nil {
		return output, err
	}

	i.Password = string(hashedPassword)

	err = as.userRepo.InsertUserByEmail(&domain.User{
		Email:    i.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return output, err
	}

	return output, nil
}

func (as *authService) Login(r *http.Request, w http.ResponseWriter, i LoginInput) error {
	// We should valid the email and password here
	valErrs, err := as.validator.Struct(i)
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	if valErrs != nil {
		return ErrInvalidCredentials
	}

	user, err := as.userRepo.SelectUserByEmail(i.Email)
	if err != nil {
		return ErrInvalidCredentials
	}

	err = as.passwordManager.CompareHashAndPassword([]byte(user.Password), i.Password)
	if err != nil {
		return ErrInvalidCredentials
	}

	session, err := as.sessionStore.New(r, SESSION_NAME)
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	// Add the user ID to the session
	session.Values[SESSION_USER_ID_FIELD] = user.ID

	// Save the session
	err = as.sessionStore.Save(r, w, session)
	if err != nil {
		return fmt.Errorf("Login: error saving the session: %w", err)
	}

	return nil
}

func (as *authService) Logout(r *http.Request, w http.ResponseWriter) error {
	session, err := as.sessionStore.Get(r, SESSION_NAME)
	if err != nil {
		return fmt.Errorf("Logout: error getting session from store %w", err)
	}

	err = as.sessionStore.Delete(r, w, session)
	if err != nil {
		return fmt.Errorf("Logout: error deleting session from store %w", err)
	}

	return nil
}

func (as *authService) ValidateRegisterEmail(i ValidateRegisterEmailInput) (ValidateRegisterEmailOutput, error) {
	output := ValidateRegisterEmailOutput{}

	valErrs, err := as.validator.Struct(i)
	if err != nil {
		return output, err
	}

	if valErrs != nil {
		output.ValidationErrors = valErrs
		return output, nil
	}

	return output, nil
}

func (as *authService) ValidateRegisterPassword(i ValidateRegisterPasswordInput) (ValidateRegisterPasswordOutput, error) {
	output := ValidateRegisterPasswordOutput{}

	valErrs, err := as.validator.Struct(i)
	if err != nil {
		return output, err
	}

	if valErrs != nil {
		output.ValidationErrors = valErrs
		return output, nil
	}

	return output, nil
}
