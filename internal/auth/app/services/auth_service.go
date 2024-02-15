package services

import (
	"net/http"

	"github.com/juanma1331/go-demo/internal/auth/domain"
	"github.com/juanma1331/go-demo/internal/shared"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

const (
	SESSION_NAME          = "go-demo-session-name"
	SESSION_USER_ID_FIELD = "user_id"
	SESSION_TABLE_NAME    = "sessions"
	SESSION_PATH          = "/"
	SESSION_MAX_AGE       = 86400
)

type AuthUserRepository interface {
	// InsertUserByEmail inserts a new user into the database using their email.
	// Returns an error in case of failure.
	InsertUserByEmail(*domain.User) error

	// SelectUserByEmail searches for a user by their email.
	// Returns ErrUserNotFound if the user is not found, and other errors in case of database problems.
	SelectUserByEmail(string) (*domain.User, error)

	// SelectUserByID searches for a user by their ID.
	// Returns ErrUserNotFound if the user is not found, and other errors in case of database problems.
	SelectUserByID(uuid.UUID) (*domain.User, error)
}

type PasswordManager interface {
	GenerateFromPassword(string) ([]byte, error)
	CompareHashAndPassword([]byte, string) error
}

type SessionStore interface {
	New(*http.Request, string) (*sessions.Session, error)
	Save(*http.Request, http.ResponseWriter, *sessions.Session) error
	Delete(*http.Request, http.ResponseWriter, *sessions.Session) error
	Get(*http.Request, string) (*sessions.Session, error)
}

type AuthService interface {
	Register(RegisterInput) (RegisterOutput, error)
	Login(*http.Request, http.ResponseWriter, LoginInput) error
	Logout(r *http.Request, w http.ResponseWriter) error
	ValidateRegisterEmail(ValidateRegisterEmailInput) (ValidateRegisterEmailOutput, error)
	ValidateRegisterPassword(ValidateRegisterPasswordInput) (ValidateRegisterPasswordOutput, error)
}

type LoginInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RegisterInput struct {
	Email    string `validate:"required,email,unique_email"`
	Password string `validate:"required"`
}

type RegisterOutput struct {
	ValidationErrors *shared.ValidationErrors
}

type ValidateRegisterEmailInput struct {
	Email string `validate:"required,email,unique_email"`
}

type ValidateRegisterEmailOutput struct {
	ValidationErrors *shared.ValidationErrors
}

type ValidateRegisterPasswordInput struct {
	Password string `validate:"min=8,max=100,required,lowercase,uppercase,number,special"`
}

type ValidateRegisterPasswordOutput struct {
	ValidationErrors *shared.ValidationErrors
}
