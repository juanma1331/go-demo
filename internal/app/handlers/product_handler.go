package handlers

import (
	"go-demo/internal/app"
	"go-demo/views/productview"

	"github.com/labstack/echo"
)

type productHandler struct{}

func NewProductHandler() *productHandler {
	return &productHandler{}
}

func (h *productHandler) HandleProductIndex(c echo.Context) error {
	cc := c.(app.AppContext)
	return cc.RenderComponent(productview.Index())
}

func (uh *productHandler) HandleProductImage(c echo.Context) error {
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
