package domain

import (
	"fmt"
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
	ID            uuid.UUID     `bun:"id,pk,type:uuid"`
	UserID        uuid.UUID     `bun:"user_id,notnull,type:uuid"`
	CreatedAt     time.Time     `bun:",nullzero,notnull,default:current_timestamp"`
	Status        string        `bun:"status,notnull"`
	CartDetails   []*CartDetail `bun:"rel:has-many,join:id=cart_id"`
}

func (c *Cart) GetCartDetail(id uuid.UUID) *CartDetail {
	for _, cd := range c.CartDetails {
		if cd.ID == id {
			return cd
		}
	}
	return nil

}

func (c *Cart) AddOrUpdateProduct(product Product) (CartDetail, bool, error) {
	found := false
	for i, detail := range c.CartDetails {
		if detail.ProductID == product.ID {
			c.CartDetails[i].Quantity += 1
			found = true
			return *c.CartDetails[i], true, nil
		}
	}

	if !found {
		newCartDetail := &CartDetail{
			ID:        uuid.New(),
			CartID:    c.ID,
			ProductID: product.ID,
			Product:   &product,
			Quantity:  1,
		}
		c.CartDetails = append(c.CartDetails, newCartDetail)
		return *newCartDetail, false, nil
	}

	return CartDetail{}, false, fmt.Errorf("error adding or updating product")
}

func (c *Cart) RemoveCartDetail(id uuid.UUID) error {
	updatedCartDetails := make([]*CartDetail, 0)
	found := false
	for _, cd := range c.CartDetails {
		if cd.ID != id {
			updatedCartDetails = append(updatedCartDetails, cd)
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("cart detail not found")
	}
	c.CartDetails = updatedCartDetails
	return nil
}

func (c *Cart) GetTotalQuantity() int {
	total := 0
	for _, cd := range c.CartDetails {
		total += cd.Quantity
	}
	return total
}

func (c *Cart) GetTotalPrice() int64 {
	total := int64(0)
	for _, cd := range c.CartDetails {
		total += cd.Product.Price * int64(cd.Quantity)
	}
	return total
}

func (c *Cart) IsEmpty() bool {
	return len(c.CartDetails) == 0
}
