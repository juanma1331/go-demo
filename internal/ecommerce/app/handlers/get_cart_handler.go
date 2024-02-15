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

type getCartHandler struct {
	db *bun.DB
}

func NewGetCartHandler(db *bun.DB) getCartHandler {
	return getCartHandler{db: db}
}

func (h getCartHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)

	carts := []domain.Cart{}
	err := h.db.NewSelect().
		Model(&carts).
		Relation("CartDetails.Product").
		Where("user_id = ?", cc.User.ID).
		Scan(cc.Request().Context())
	if err != nil {
		return fmt.Errorf("HandleGetCart: error getting cart: %w", err)
	}

	token := csrf.Token(cc.Request())

	if len(carts) == 0 {
		cart := domain.Cart{
			ID:     uuid.New(),
			Status: domain.CART_STATUS_ACTIVE,
			UserID: cc.User.ID,
		}

		_, err := h.db.NewInsert().Model(&cart).Exec(cc.Request().Context())
		if err != nil {
			return err
		}

		return cc.RenderComponent(layouts.Cart([]layouts.CartProductViewModel{}, token))
	}

	cart := carts[0]

	if len(cart.CartDetails) == 0 {
		return cc.RenderComponent(layouts.Cart([]layouts.CartProductViewModel{}, token))
	}

	cartProducts := make([]layouts.CartProductViewModel, 0, len(cart.CartDetails))
	totalProductsQuantity := 0
	for _, cartDetail := range cart.CartDetails {
		cartProducts = append(cartProducts, layouts.CartProductViewModel{
			DetailID:           cartDetail.ID.String(),
			ProductID:          cartDetail.ProductID.String(),
			ProductName:        cartDetail.Product.Name,
			ProductDescription: cartDetail.Product.Description,
			ProductPrice:       cartDetail.Product.Price,
			Quantity:           cartDetail.Quantity,
		})
		totalProductsQuantity += cartDetail.Quantity
	}

	cartUpdatedTrigger := shared.HtmxTrigger{
		Name: "cart_updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", totalProductsQuantity),
		},
	}

	err = shared.SetHtmxTriggers(cc.Response().Writer, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleGetCart: error setting htmx trigger: %w", err)
	}

	return cc.RenderComponent(layouts.Cart(cartProducts, token))
}
