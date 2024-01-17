package domain

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AuthToken struct {
	bun.BaseModel `bun:"auth_tokens"`

	ID     uuid.UUID `bun:"id,pk,type:uuid"`
	Token  string    `bun:"token,notnull"`
	UserID uuid.UUID `bun:"user_id,notnull"`
}
