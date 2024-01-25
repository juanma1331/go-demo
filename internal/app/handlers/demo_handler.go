package handlers

import (
	"database/sql"
	"fmt"
	"go-demo/internal/app"
	"go-demo/internal/domain"
	"go-demo/views/demoview"
	"go-demo/views/shared"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type productHandler struct {
	db *bun.DB
}

func NewDemoHandler(db *bun.DB) *productHandler {
	return &productHandler{
		db: db,
	}
}

func (h *productHandler) HandleProductIndex(c echo.Context) error {
	cc := c.(app.AppContext)
	return cc.RenderComponent(demoview.IndexPage())
}

func (h *productHandler) GetProductList(c echo.Context) error {
	cc := c.(app.AppContext)

	products := []domain.Product{}
	h.db.NewSelect().Model(&products).Scan(cc.Request().Context())

	return cc.RenderComponent(demoview.ProductList(products))
}

func (h *productHandler) HandleGetCart(c echo.Context) error {
	cc := c.(app.AppContext)

	carts := []domain.Cart{}
	err := h.db.NewSelect().
		Model(&carts).
		Relation("CartDetails.Product").
		Where("user_id = ?", cc.User.ID).
		Scan(cc.Request().Context())
	if err != nil {
		return fmt.Errorf("HandleGetCart: error getting cart: %w", err)
	}

	fmt.Printf("Cart: %v\n", carts)

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

		return cc.RenderComponent(demoview.Cart(cart.ID.String(), []domain.Product{}))
	}

	cart := carts[0]

	if len(cart.CartDetails) == 0 {
		return cc.RenderComponent(demoview.Cart(cart.ID.String(), []domain.Product{}))
	}

	products := make([]domain.Product, 0, len(cart.CartDetails))
	for _, cartDetail := range cart.CartDetails {
		products = append(products, *cartDetail.Product)
	}

	fmt.Printf("cart is: %v\n", cart)

	return cc.RenderComponent(demoview.Cart(cart.ID.String(), products))
}

func (h *productHandler) HandleAddToCart(c echo.Context) error {
	cc := c.(app.AppContext)

	userId := cc.User.ID
	productId := c.FormValue("product_id")

	// Check if product exists
	var product domain.Product
	err := h.db.
		NewSelect().
		Model(&product).
		Where("id = ?", productId).
		Scan(cc.Request().
			Context())
	if err != nil {
		return fmt.Errorf("HandleAddToCart: product not found: %w", err)
	}

	// Check if cart exists
	var cart domain.Cart
	err = h.db.
		NewSelect().
		Model(&cart).
		Where("user_id = ? AND status = ?", userId, domain.CART_STATUS_ACTIVE).
		Scan(cc.Request().Context())

	if err != nil {
		return fmt.Errorf("HandleAddToCart: cart not found: %w", err)
	}

	// Check if product is already in cart
	var cartDetail domain.CartDetail
	err = h.db.
		NewSelect().
		Model(&cartDetail).
		Where("cart_id = ? AND product_id = ?", cart.ID, productId).
		Scan(cc.Request().Context())

	if err == nil {
		// Product is already in cart increment quantity
		cartDetail.Quantity++
		_, err := h.db.NewUpdate().Model(&cartDetail).Exec(cc.Request().Context())
		if err != nil {
			return fmt.Errorf("HandleAddToCart: error updating cart detail: %w", err)
		}

	} else {
		if err != sql.ErrNoRows {
			return fmt.Errorf("HandleAddToCart: error checking cart detail: %w", err)
		}
		// Product is not in cart, add it

		cartDetail = domain.CartDetail{
			ID:        uuid.New(),
			CartID:    cart.ID,
			ProductID: product.ID,
			Quantity:  1,
		}

		_, err := h.db.NewInsert().Model(&cartDetail).Exec(cc.Request().Context())
		if err != nil {
			return fmt.Errorf("HandleAddToCart: error inserting cart detail: %w", err)
		}
	}

	// This will trigger an event on the client side
	// Pretty useful for notifications like toasts, alerts, etc.
	trigger := app.HtmxTrigger{
		Name:  shared.NOTIFY_TRIGGER_NAME,
		Value: map[string]string{"message": "Product Added", "type": "success"},
	}

	if err := app.SetHtmxTrigger(c.Response(), trigger); err != nil {
		return err
	}

	return cc.RenderComponent(demoview.CartProduct(product))
}

func (uh *productHandler) HandleProductImage(c echo.Context) error {
	productId := c.Param("id")
	imageSize := c.Param("size")

	if imageSize != "small" && imageSize != "medium" {
		return c.String(http.StatusBadRequest, "Invalid image size")
	}

	var product domain.Product
	err := uh.db.NewSelect().Model(&product).Where("id = ?", productId).Scan(c.Request().Context())
	if err != nil {
		return c.String(http.StatusNotFound, "Product Not Found")
	}

	var productImage []byte
	switch imageSize {
	case "small":
		productImage = product.ImageSmall
	case "medium":
		productImage = product.ImageMedium
	}

	return c.Blob(http.StatusOK, "image/jpeg", productImage)
}
