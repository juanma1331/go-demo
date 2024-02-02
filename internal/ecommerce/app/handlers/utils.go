package handlers

import "go-demo/internal/ecommerce/domain"

func calculateTotalQuantity(cartDetails []domain.CartDetail) int {
	total := 0
	for _, cd := range cartDetails {
		total += cd.Quantity
	}
	return total
}
