package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	CART_STATUS_ACTIVE = "active"
	CART_STATUS_PAID   = "paid"
)

type Cart struct {
	bun.BaseModel `bun:"carts,alias:c"`
	ID            uuid.UUID    `bun:"id,pk,type:uuid"`
	UserID        uuid.UUID    `bun:"user_id,notnull"`
	CreationDate  time.Time    `bun:"creation_date,notnull"`
	Status        string       `bun:"status,notnull"`
	CartDetails   []CartDetail `bun:"rel:has-many,join:id=cart_id"`
}
