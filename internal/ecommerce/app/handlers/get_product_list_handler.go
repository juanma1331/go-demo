package handlers

import (
	"go-demo/internal/ecommerce/domain"
	"go-demo/internal/shared"
	"go-demo/views/ecommerce"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type getProductListHandler struct {
	db *bun.DB
}

func NewGetProductListHandler(db *bun.DB) getProductListHandler {
	return getProductListHandler{db: db}
}

func (h getProductListHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)

	products := []domain.Product{}
	h.db.NewSelect().Model(&products).Scan(cc.Request().Context())

	csrfToken := csrf.Token(cc.Request())

	return cc.RenderComponent(ecommerce.ProductList(products, csrfToken))
}
