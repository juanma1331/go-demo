package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"products"`

	ID          uuid.UUID `bun:"id,pk,type:uuid"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description,notnull"`
	Price       int64     `bun:"price,notnull"`
	Image       []byte    `bun:"type:blob"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
