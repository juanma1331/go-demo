package handlers

import (
	"go-demo/internal/ecommerce/domain"
	"net/http"

	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type getProductImageHandler struct {
	db *bun.DB
}

func NewGetProductImageHandler(db *bun.DB) getProductImageHandler {
	return getProductImageHandler{db: db}
}

func (h getProductImageHandler) Handler(c echo.Context) error {
	productId := c.Param("id")
	imageSize := c.Param("size")

	if imageSize != "small" && imageSize != "medium" {
		return c.String(http.StatusBadRequest, "Invalid image size")
	}

	var product domain.Product
	err := h.db.NewSelect().Model(&product).Where("id = ?", productId).Scan(c.Request().Context())
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
