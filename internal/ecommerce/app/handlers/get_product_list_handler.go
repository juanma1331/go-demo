package handlers

import (
	"net/http"

	"github.com/juanma1331/go-demo/internal/shared"
	"github.com/juanma1331/go-demo/views/ecommerce"

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

	limit := 30
	var initialCursor int64

	products, newCursor, err := selectProductsNextPage(c, h.db, initialCursor, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	productViewModels := []ecommerce.ProductViewModel{}
	for _, p := range products {
		productViewModels = append(productViewModels, ecommerce.ProductViewModel{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	csrfToken := csrf.Token(cc.Request())

	return cc.RenderComponent(ecommerce.ProductList(productViewModels, csrfToken, newCursor))
}
