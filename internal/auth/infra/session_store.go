package infra

import (
	"encoding/gob"
	"fmt"
	"go-demo/internal/shared"
	"net/http"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type SessionStore struct {
	pgstore.PGStore
}

func init() {
	gob.Register(uuid.UUID{})
	gob.Register(time.Time{})
	gob.Register(shared.FlashMessage{})
}

func NewSessionStore(dsn string) (*SessionStore, error) {
	store, err := pgstore.NewPGStore(dsn, []byte("your-key"))
	if err != nil {
		return nil, err
	}

	return &SessionStore{
		PGStore: *store,
	}, nil
}

func (*SessionStore) Delete(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	s.Options.MaxAge = -1
	if err := s.Save(r, w); err != nil {
		return fmt.Errorf("Delete: error saving session: %w", err)
	}

	return nil
}
