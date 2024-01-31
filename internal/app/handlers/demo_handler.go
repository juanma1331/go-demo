package handlers

import (
	"fmt"
	"go-demo/internal/app"
	"go-demo/internal/domain"
	"go-demo/views/demoview"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type demoHandler struct {
	db *bun.DB
}

func NewDemoHandler(db *bun.DB) *demoHandler {
	return &demoHandler{
		db: db,
	}
}

func (h *demoHandler) HandleProductIndex(c echo.Context) error {
	cc := c.(app.AppContext)
	return cc.RenderComponent(demoview.IndexPage())
}

func (h *demoHandler) GetProductList(c echo.Context) error {
	cc := c.(app.AppContext)

	products := []domain.Product{}
	h.db.NewSelect().Model(&products).Scan(cc.Request().Context())

	csrfToken := csrf.Token(cc.Request())

	return cc.RenderComponent(demoview.ProductList(products, csrfToken))
}

func (h *demoHandler) HandleGetCart(c echo.Context) error {
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

		return cc.RenderComponent(demoview.Cart([]demoview.CartProductViewModel{}, token))
	}

	cart := carts[0]

	if len(cart.CartDetails) == 0 {
		return cc.RenderComponent(demoview.Cart([]demoview.CartProductViewModel{}, token))
	}

	cartProducts := make([]demoview.CartProductViewModel, 0, len(cart.CartDetails))
	totalProductsQuantity := 0
	for _, cartDetail := range cart.CartDetails {
		cartProducts = append(cartProducts, demoview.CartProductViewModel{
			DetailId:           cartDetail.ID.String(),
			ProductName:        cartDetail.Product.Name,
			ProductDescription: cartDetail.Product.Description,
			Quantity:           cartDetail.Quantity,
		})
		totalProductsQuantity += cartDetail.Quantity
	}

	cartUpdatedTrigger := app.HtmxTrigger{
		Name: "cart-updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", totalProductsQuantity),
		},
	}

	err = app.SetHtmxTriggers(cc.Response().Writer, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleGetCart: error setting htmx trigger: %w", err)
	}

	return cc.RenderComponent(demoview.Cart(cartProducts, token))
}

func (h *demoHandler) HandleAddToCart(c echo.Context) error {
	cc := c.(app.AppContext)
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
		app.SetHtmxRetarget(cc.Response().Writer, fmt.Sprintf("#cart-item-%s", cartDetail.ID.String()))

		app.SetHtmxReswap(cc.Response().Writer, "outerHTML")
	}

	// Notifying the client that the cart was updated
	notifyTrigger := app.HtmxTrigger{
		Name: "notify",
		Value: map[string]string{
			"message": fmt.Sprintf("Product with id=%s has been added", productId),
			"type":    "success",
		},
	}

	cartUpdatedTrigger := app.HtmxTrigger{
		Name: "cart-updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", calculateTotalQuantity(cart.CartDetails)),
		},
	}

	err = app.SetHtmxTriggers(cc.Response().Writer, notifyTrigger, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleAddToCart: error setting htmx triggers: %w", err)
	}

	// If the product was not in the cart, we need to add it to the cart
	return cc.RenderComponent(demoview.CartProduct(demoview.CartProductViewModel{
		DetailId:           cartDetail.ID.String(),
		ProductName:        product.Name,
		ProductDescription: product.Description,
		Quantity:           cartDetail.Quantity,
	}, token))
}

func (h *demoHandler) getProduct(productId string, cc app.AppContext) (domain.Product, error) {
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

func (h *demoHandler) getActiveCart(userId string, cc app.AppContext) (domain.Cart, error) {
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

func (h *demoHandler) updateOrInsertProductInCart(cart *domain.Cart, product domain.Product, cc app.AppContext) (domain.CartDetail, error) {
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

func calculateTotalQuantity(cartDetails []domain.CartDetail) int {
	total := 0
	for _, cd := range cartDetails {
		total += cd.Quantity
	}
	return total
}

func (h *demoHandler) HandleDecreaseQuantity(c echo.Context) error {
	cc := c.(app.AppContext)
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

	// Decrease quantity
	cartDetail.Quantity--

	_, err = h.db.NewUpdate().
		Model(cartDetail).
		Where("id = ?", cartDetail.ID).
		Column("quantity").
		Exec(cc.Request().Context())
	if err != nil {
		return fmt.Errorf("HandleDecreaseQuantity: error updating cart detail: %w", err)
	}

	// Notifying the client that the cart was updated
	notifyTrigger := app.HtmxTrigger{
		Name: "notify",
		Value: map[string]string{
			"message": fmt.Sprintf("Product with id=%s has been removed", cartDetail.ProductID.String()),
			"type":    "success",
		},
	}

	cartUpdatedTrigger := app.HtmxTrigger{
		Name: "cart-updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", calculateTotalQuantity(cart.CartDetails)),
		},
	}

	fmt.Printf("total quantity: %d\n", calculateTotalQuantity(cart.CartDetails))

	err = app.SetHtmxTriggers(cc.Response().Writer, notifyTrigger, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleDecreaseQuantity: error setting htmx triggers: %w", err)
	}

	return cc.RenderComponent(demoview.CartProduct(demoview.CartProductViewModel{
		DetailId:           cartDetail.ID.String(),
		ProductName:        cartDetail.Product.Name,
		ProductDescription: cartDetail.Product.Description,
		Quantity:           cartDetail.Quantity,
	}, token))
}

func (h *demoHandler) HandleRemoveFromCart(c echo.Context) error {
	cc := c.(app.AppContext)
	cartDetailId := c.FormValue("cart_detail_id")

	// Get user's active cart
	cart, err := h.getActiveCart(cc.User.ID.String(), cc)
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error getting active cart: %w", err)
	}

	fmt.Printf("cart id: %s\n", cartDetailId)

	// Get cart detail
	var cartDetail domain.CartDetail
	for _, cd := range cart.CartDetails {
		if cd.ID.String() == cartDetailId {
			cartDetail = cd
			break
		}
	}

	if cartDetail.ID == uuid.Nil {
		return fmt.Errorf("HandleRemoveFromCart: cart detail not found")
	}

	if cartDetail.Quantity > 1 {
		return fmt.Errorf("HandleRemoveFromCart: cart detail quantity is more than 1")
	}

	// Remove cart detail
	_, err = h.db.NewDelete().
		Model(&cartDetail).
		Where("id = ?", cartDetail.ID).
		Exec(cc.Request().Context())
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error deleting cart detail: %w", err)
	}

	// Notifying the client that the cart was updated
	notifyTrigger := app.HtmxTrigger{
		Name: "notify",
		Value: map[string]string{
			"message": fmt.Sprintf("Product with id=%s has been removed", cartDetail.ProductID.String()),
			"type":    "success",
		},
	}

	cartUpdatedTrigger := app.HtmxTrigger{
		Name: "cart-updated",
		Value: map[string]string{
			"quantity": fmt.Sprintf("%d", calculateTotalQuantity(cart.CartDetails)-1),
		},
	}

	err = app.SetHtmxTriggers(cc.Response().Writer, notifyTrigger, cartUpdatedTrigger)
	if err != nil {
		return fmt.Errorf("HandleRemoveFromCart: error setting htmx triggers: %w", err)
	}

	return cc.NoContent(http.StatusOK)
}

func (uh *demoHandler) HandleProductImage(c echo.Context) error {
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
