package handlers

import (
	"fmt"
	"net/http"

	"github.com/juanma1331/go-demo/internal/ecommerce/domain"
	"github.com/juanma1331/go-demo/internal/shared"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type removeFromCartHandler struct {
	db *bun.DB
}

func NewRemoveFromCartHandler(db *bun.DB) removeFromCartHandler {
	return removeFromCartHandler{
		db: db,
	}
}

func (h removeFromCartHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)
	cartDetailId := c.FormValue("cart_detail_id")

	// Get user's active cart
	cart, err := h.getActiveCart(cc.User.ID.String(), cc)
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error getting active cart: %w", err)
	}

	// Get cart detail
	cartDetail := cart.GetCartDetail(uuid.MustParse(cartDetailId))

	if cartDetail.ID == uuid.Nil {
		return fmt.Errorf("HandleRemoveFromCart: cart detail not found")
	}

	if cartDetail.Quantity > 1 {
		return fmt.Errorf("HandleRemoveFromCart: cart detail quantity is more than 1")
	}

	// Remove cart detail from DB
	err = h.removeCartDetailInDB(*cartDetail, cc)
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error removing cart detail in db: %w", err)
	}

	// Update cart in memory
	err = cart.RemoveCartDetail(cartDetail.ID)
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error removing cart detail: %w", err)
	}

	// Notifying the client that the cart was updated
	notifyTrigger := shared.HtmxTrigger{
		Name: "notify",
		Value: map[string]string{
			"message": fmt.Sprintf("A product has been removed from your cart: %s", cartDetail.Product.Name),
			"type":    "success",
		},
	}

	cartUpdatedTrigger := shared.HtmxTrigger{
		Name: "cart_updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", cart.GetTotalQuantity()),
			"total":    fmt.Sprintf("%d", cart.GetTotalPrice()),
		},
	}

	err = shared.SetHtmxTriggers(cc.Response().Writer, notifyTrigger, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error setting htmx triggers: %w", err)
	}

	return cc.NoContent(http.StatusOK)
}

func (h removeFromCartHandler) getActiveCart(userId string, cc shared.AppContext) (domain.Cart, error) {
	var cart domain.Cart
	err := h.db.
		NewSelect().
		Model(&cart).
		Relation("CartDetails.Product").
		Where("user_id = ? AND status = ?", userId, domain.CART_STATUS_ACTIVE).
		Scan(cc.Request().Context())

	if err != nil {
		return domain.Cart{}, fmt.Errorf("getActiveCart: error getting cart: %w", err)
	}
	return cart, nil
}

func (h removeFromCartHandler) removeCartDetailInDB(cartDetail domain.CartDetail, cc shared.AppContext) error {
	_, err := h.db.NewDelete().
		Model(&cartDetail).
		Where("id = ?", cartDetail.ID).
		Exec(cc.Request().Context())
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error deleting cart detail: %w", err)
	}

	return nil
}
