package handlers

import (
	"go-demo/internal/app"
	"go-demo/view/productview"

	"github.com/labstack/echo"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) HandleProductIndex(c echo.Context) error {
	cc := c.(app.AppContext)
	return cc.RenderComponent(productview.Index())
}

func (uh *ProductHandler) HandleProductImage(c echo.Context) error {
	// productId, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "Invalid product ID")
	// }

	// var product domain.Product
	// err = uh.DB.NewSelect().Model(&product).Where("id = ?", productId).Scan(c.Request().Context())
	// if err != nil {
	// 	return c.String(http.StatusNotFound, "Product Not Found")
	// }

	// return c.Blob(http.StatusOK, "image/jpeg", product.Image)

	return nil
}
