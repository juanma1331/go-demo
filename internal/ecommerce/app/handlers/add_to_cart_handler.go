package handlers

import (
	"fmt"
	"go-demo/internal/ecommerce/domain"
	"go-demo/internal/shared"
	"go-demo/views/layouts"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type addToCartHandler struct {
	db *bun.DB
}

func NewAddToCartHandler(db *bun.DB) addToCartHandler {
	return addToCartHandler{db: db}
}

func (h addToCartHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)
	productId := c.FormValue("product_id")
	token := csrf.Token(cc.Request())

	product, err := h.getProduct(productId, cc)
	if err != nil {
		return fmt.Errorf("HandleAddToCart: product not found: %w", err)
	}

	cart, err := h.getActiveCart(cc.User.ID.String(), cc)
	if err != nil {
		return fmt.Errorf("HandleAddToCart: cart not found: %w", err)
	}

	cartDetail, err := h.updateOrInsertProductInCart(&cart, product, cc)
	if err != nil {
		return fmt.Errorf("HandleAddToCart: error updating or inserting product in cart: %w", err)
	}

	if cartDetail.Quantity > 1 {
		// If the product was already in the cart, we need to update the quantity on the client cart
		shared.SetHtmxRetarget(cc.Response().Writer, fmt.Sprintf("#cart-item-%s", cartDetail.ID.String()))

		shared.SetHtmxReswap(cc.Response().Writer, "outerHTML")
	}

	// Notifying the client that the cart was updated
	notifyTrigger := shared.HtmxTrigger{
		Name: "notify",
		Value: map[string]string{
			"message": fmt.Sprintf("Product with id=%s has been added", productId),
			"type":    "success",
		},
	}

	cartUpdatedTrigger := shared.HtmxTrigger{
		Name: "cart_updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", calculateTotalQuantity(cart.CartDetails)),
		},
	}

	err = shared.SetHtmxTriggers(cc.Response().Writer, notifyTrigger, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleAddToCart: error setting htmx triggers: %w", err)
	}

	// If the product was not in the cart, we need to add it to the cart
	return cc.RenderComponent(layouts.CartProduct(layouts.CartProductViewModel{
		DetailID:           cartDetail.ID.String(),
		ProductID:          product.ID.String(),
		ProductName:        product.Name,
		ProductDescription: product.Description,
		Quantity:           cartDetail.Quantity,
	}, token))
}

func (h addToCartHandler) getProduct(productId string, cc shared.AppContext) (domain.Product, error) {
	var product domain.Product
	err := h.db.
		NewSelect().
		Model(&product).
		Where("id = ?", productId).
		Scan(cc.Request().Context())
	if err != nil {
		return domain.Product{}, fmt.Errorf("getProduct: error getting product: %w", err)
	}
	return product, nil
}

func (h addToCartHandler) getActiveCart(userId string, cc shared.AppContext) (domain.Cart, error) {
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

func (h addToCartHandler) updateOrInsertProductInCart(cart *domain.Cart, product domain.Product, cc shared.AppContext) (domain.CartDetail, error) {
	var cartDetail domain.CartDetail
	productInCart := false

	for i, cd := range cart.CartDetails {
		if cd.ProductID == product.ID {
			productInCart = true
			cart.CartDetails[i].Quantity++
			cartDetail = cart.CartDetails[i]
			_, err := h.db.NewUpdate().
				Model(&cartDetail).
				Where("id = ?", cartDetail.ID).
				Column("quantity").
				Exec(cc.Request().Context())
			if err != nil {
				return domain.CartDetail{}, fmt.Errorf("updateOrInsertProductInCart: error updating cart detail: %w", err)
			}
			break
		}
	}

	if !productInCart {
		cartDetail = domain.CartDetail{
			ID:        uuid.New(),
			CartID:    cart.ID,
			ProductID: product.ID,
			Quantity:  1,
		}
		cart.CartDetails = append(cart.CartDetails, cartDetail)
		_, err := h.db.NewInsert().Model(&cartDetail).Exec(cc.Request().Context())
		if err != nil {
			return domain.CartDetail{}, fmt.Errorf("updateOrInsertProductInCart: error inserting cart detail: %w", err)
		}
	}

	return cartDetail, nil
}
