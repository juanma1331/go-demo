package shared

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

const (
	FLASH_SESSION_NAME = "flash-session"
	FLASH_KEY          = "flash-messages"
)

type FlashStore interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
	Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error
	LoadFlash(c echo.Context) (*[]FlashMessage, error)
}

type flashStore struct {
	*sessions.CookieStore
}

func NewFlashStore(secret string) FlashStore {
	return &flashStore{
		CookieStore: sessions.NewCookieStore([]byte(secret)),
	}
}

func (fs flashStore) LoadFlash(c echo.Context) (*[]FlashMessage, error) {
	session, err := getSession(fs, c.Request())
	if err != nil {
		return nil, fmt.Errorf("LoadFlash: Failed to get session: %w", err)
	}

	var flashMessages []FlashMessage
	flashes := session.Flashes(FLASH_KEY)
	for _, flash := range flashes {
		if fm, ok := flash.(FlashMessage); ok {
			flashMessages = append(flashMessages, fm)
		}

	}

	session.Save(c.Request(), c.Response())

	if len(flashMessages) == 0 {
		return nil, nil
	}

	return &flashMessages, nil
}

type FlashMessage struct {
	Message string
	Type    string
}

func NewFlashMessage(message string, messageType string) *FlashMessage {
	return &FlashMessage{
		Message: message,
		Type:    messageType,
	}
}

func (fm *FlashMessage) AddToSession(store FlashStore, r *http.Request, w http.ResponseWriter) error {
	session, err := getSession(store, r)
	if err != nil {
		return fmt.Errorf("AddToSession: Failed to get session: %w", err)
	}

	session.AddFlash(fm, FLASH_KEY)

	return session.Save(r, w)
}

func getSession(store FlashStore, r *http.Request) (*sessions.Session, error) {
	return store.Get(r, FLASH_SESSION_NAME)
}
