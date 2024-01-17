package services

import (
	"go-demo/internal/domain"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var (
	SESSION_NAME        = "go-demo-session-name"
	SESSION_TOKEN_FIELD = "token"
)

type UserRepository interface {
	// InsertUserByEmail inserts a new user into the database using their email.
	// Returns an error in case of failure.
	InsertUserByEmail(*domain.User) error

	// SelectUserByEmail searches for a user by their email.
	// Returns ErrUserNotFound if the user is not found, and other errors in case of database problems.
	SelectUserByEmail(string) (*domain.User, error)

	// SelectUserByID searches for a user by their ID.
	// Returns ErrUserNotFound if the user is not found, and other errors in case of database problems.
	SelectUserByID(string) (*domain.User, error)
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

type AuthTokenManager interface {
	GenerateToken() uuid.UUID
	AddTokenToSession(*sessions.Session, uuid.UUID)
	ExtractTokenFromSession(*sessions.Session) (*uuid.UUID, error)
}

type AuthTokenRepository interface {
	// InsertToken inserts a new token into the database.
	InsertToken(*domain.AuthToken) error

	// SelectToken searches for a token by its value.
	// Returns ErrTokenNotFound if the token is not found, and other errors in case of database problems.
	SelectToken(token string) (domain.AuthToken, error)

	// DeleteToken deletes a token from the database.
	// Returns ErrTokenNotFound if the token is not found, and other errors in case of database problems.
	DeleteToken(token string) error
}

type AuthService interface {
	Register(RegisterInput) (RegisterOutput, error)
	Login(*http.Request, http.ResponseWriter, LoginInput) error
	Logout(r *http.Request, w http.ResponseWriter) error
}

type LoginInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RegisterInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RegisterOutput struct {
	ValidationErrors []ValidationError
}
