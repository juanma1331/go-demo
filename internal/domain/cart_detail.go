package domain

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CartDetail struct {
	bun.BaseModel `bun:"cart_details,alias:cd"`
	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	CartID        uuid.UUID `bun:"cart_id,notnull"`
	ProductID     uuid.UUID `bun:"product_id,notnull"`
	Product       *Product  `bun:"rel:belongs-to,join:product_id=id"`
	Quantity      int       `bun:"quantity,notnull"`
}
