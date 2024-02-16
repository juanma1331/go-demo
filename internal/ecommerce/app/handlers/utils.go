package handlers

import (
	"fmt"

	"github.com/juanma1331/go-demo/internal/ecommerce/domain"
)

func calculateTotalQuantity(cartDetails []domain.CartDetail) int {
	total := 0
	for _, cd := range cartDetails {
		total += cd.Quantity
	}
	return total
}

func calculateTotalPrice(cartDetails []domain.CartDetail) int64 {
	var total int64 = 0
	for _, cd := range cartDetails {
		// check if the product is nil
		if cd.Product == nil {
			fmt.Println("Product is nil")
		}
		total += cd.Product.Price * int64(cd.Quantity)
	}
	return total
}
