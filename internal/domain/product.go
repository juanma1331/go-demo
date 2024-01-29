package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"products,alias:p"`

	ID          uuid.UUID `bun:"id,pk,type:uuid"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description,notnull"`
	Price       int64     `bun:"price,notnull"`
	ImageSmall  []byte    `bun:"type:bytea,notnull"`
	ImageMedium []byte    `bun:"type:bytea,notnull"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
