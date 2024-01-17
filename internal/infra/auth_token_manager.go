package infra

import (
	"errors"
	"go-demo/internal/app/services"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type authTokenManager struct{}

func NewAuthTokenManager() *authTokenManager {
	return &authTokenManager{}
}

func (atm *authTokenManager) GenerateToken() uuid.UUID {
	return uuid.New()
}

func (as *authTokenManager) AddTokenToSession(s *sessions.Session, t uuid.UUID) {
	s.Values[services.SESSION_TOKEN_FIELD] = t
}

func (as *authTokenManager) ExtractTokenFromSession(s *sessions.Session) (*uuid.UUID, error) {
	token, ok := s.Values[services.SESSION_TOKEN_FIELD].(uuid.UUID)
	if !ok {
		return nil, errors.New("ExtractTokenFromSession: token not found in session")
	}
	return &token, nil
}
