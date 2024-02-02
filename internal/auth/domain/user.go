package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"users,alias:u"`
	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	Email         string    `bun:"email,unique,notnull"`
	Password      string    `bun:"password,notnull"`
	IsAdmin       bool      `bun:"is_admin,notnull,default:false"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
