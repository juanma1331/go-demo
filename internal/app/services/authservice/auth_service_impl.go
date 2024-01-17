package authservice

import (
	"fmt"
	"go-demo/internal/app/services"
	"go-demo/internal/domain"
	"net/http"
)

type AuthServiceParams struct {
	UserRepository   AuthUserRepository
	AuthTokenRepo    AuthTokenRepository
	PasswordManager  PasswordManager
	AuthTokenManager AuthTokenManager
	SessionStore     SessionStore
	Validator        services.Validator
}

type authService struct {
	userRepo         AuthUserRepository
	authTokenRepo    AuthTokenRepository
	passwordManager  PasswordManager
	authTokenManager AuthTokenManager
	sessionStore     SessionStore
	validator        services.Validator
}

func NewAuthService(params AuthServiceParams) *authService {
	return &authService{
		userRepo:         params.UserRepository,
		authTokenRepo:    params.AuthTokenRepo,
		passwordManager:  params.PasswordManager,
		authTokenManager: params.AuthTokenManager,
		sessionStore:     params.SessionStore,
		validator:        params.Validator,
	}
}

func (as *authService) Register(i RegisterInput) (RegisterOutput, error) {
	output := RegisterOutput{}

	valErrs, err := as.validator.Struct(i)
	if err != nil {
		return output, err
	}

	if len(valErrs) > 0 {
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

	if len(valErrs) > 0 {
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

	// Generate a new token
	token := as.authTokenManager.GenerateToken()

	// Put the token in the session
	as.authTokenManager.AddTokenToSession(session, token)

	// Save token in the database
	err = as.authTokenRepo.InsertToken(&domain.AuthToken{
		Token:  token.String(),
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	// Save the session
	err = as.sessionStore.Save(r, w, session)
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	return nil
}

func (as *authService) Logout(r *http.Request, w http.ResponseWriter) error {
	session, err := as.sessionStore.Get(r, SESSION_NAME)
	if err != nil {
		return err
	}

	token, err := as.authTokenManager.ExtractTokenFromSession(session)
	if err != nil {
		return err
	}

	err = as.authTokenRepo.DeleteToken(token.String())
	if err != nil {
		return err
	}

	return as.sessionStore.Delete(r, w, session)
}
