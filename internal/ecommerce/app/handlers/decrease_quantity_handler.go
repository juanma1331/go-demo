package handlers

import (
	"fmt"

	"github.com/juanma1331/go-demo/internal/ecommerce/domain"
	"github.com/juanma1331/go-demo/internal/shared"
	"github.com/juanma1331/go-demo/views/layouts"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type decreaseQuantityHandler struct {
	db *bun.DB
}

func NewDecreaseQuantityHandler(db *bun.DB) decreaseQuantityHandler {
	return decreaseQuantityHandler{db: db}
}

func (h decreaseQuantityHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)
	cartDetailId := c.FormValue("cart_detail_id")
	token := csrf.Token(cc.Request())

	// Get user's active cart
	cart, err := h.getActiveCart(cc.User.ID.String(), cc)
	if err != nil {
		return fmt.Errorf("HandleDecreaseQuantity: error getting active cart: %w", err)
	}

	// Get cart detail
	var cartDetailIndex int
	for i, cd := range cart.CartDetails {
		if cd.ID.String() == cartDetailId {
			cartDetailIndex = i
			break
		}
	}

	cartDetail := &cart.CartDetails[cartDetailIndex]

	if cartDetail.ID == uuid.Nil {
		return fmt.Errorf("HandleDecreaseQuantity: cart detail not found")
	}

	if cartDetail.Quantity == 1 {
		return fmt.Errorf("HandleDecreaseQuantity: cart detail quantity is already 1")
	}

	cartDetail.Quantity--

	err = h.updateCartDetail(cartDetail, cc)
	if err != nil {
		return fmt.Errorf("HandleDecreaseQuantity: error updating cart detail: %w", err)
	}

	// Notifying the client that the cart was updated
	notifyTrigger := shared.HtmxTrigger{
		Name: "notify",
		Value: map[string]string{
			"message": fmt.Sprintf("Product with id=%s has been removed", cartDetail.ProductID.String()),
			"type":    "success",
		},
	}

	// Notifying the client that the cart was updated
	cartUpdatedTrigger := shared.HtmxTrigger{
		Name: "cart_updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", calculateTotalQuantity(cart.CartDetails)),
		},
	}

	err = shared.SetHtmxTriggers(cc.Response().Writer, notifyTrigger, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleDecreaseQuantity: error setting htmx triggers: %w", err)
	}

	return cc.RenderComponent(layouts.CartProduct(layouts.CartProductViewModel{
		DetailID:           cartDetail.ID.String(),
		ProductID:          cartDetail.ProductID.String(),
		ProductName:        cartDetail.Product.Name,
		ProductDescription: cartDetail.Product.Description,
		ProductPrice:       cartDetail.Product.Price * int64(cartDetail.Quantity),
		Quantity:           cartDetail.Quantity,
	}, token))
}

func (h decreaseQuantityHandler) getActiveCart(userId string, cc shared.AppContext) (domain.Cart, error) {
	var cart domain.Cart
	err := h.db.
		NewSelect().
		Model(&cart).
		Relation("CartDetails.Product").
		Where("user_id = ? AND status = ?", userId, domain.CART_STATUS_ACTIVE).
		Scan(cc.Request().Context())

	if err != nil {
		return domain.Cart{}, err
	}
	return cart, nil
}

func (h decreaseQuantityHandler) updateCartDetail(
	cartDetail *domain.CartDetail,
	cc shared.AppContext,
) error {
	_, err := h.db.NewUpdate().
		Model(cartDetail).
		Where("id = ?", cartDetail.ID).
		Column("quantity").
		Exec(cc.Request().Context())

	return err
}
