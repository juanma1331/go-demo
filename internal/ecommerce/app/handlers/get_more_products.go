package handlers

import (
	"fmt"
	"go-demo/internal/ecommerce/domain"
	"go-demo/internal/shared"
	"go-demo/views/ecommerce"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type getMoreProductsHandler struct {
	db *bun.DB
}

func NewGetMoreProductsHandler(db *bun.DB) getMoreProductsHandler {
	return getMoreProductsHandler{db: db}
}

func (h getMoreProductsHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)

	limit := 30
	cursor, err := strconv.ParseInt(cc.Param("cursor"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	products, newCursor, err := selectProductsNextPage(c, h.db, cursor, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	productViewModels := []ecommerce.ProductViewModel{}
	for _, p := range products {
		productViewModels = append(productViewModels, ecommerce.ProductViewModel{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
		})
	}

	csrfToken := csrf.Token(cc.Request())

	return cc.RenderComponent(ecommerce.ProductList(productViewModels, csrfToken, newCursor))
}

func selectProductsNextPage(ctx echo.Context, db *bun.DB, cursor int64, limit int) ([]domain.Product, string, error) {
	var products []domain.Product
	query := db.NewSelect().Model(&products).Limit(limit + 1) // We are asking for one more to check if we have reached the end

	if cursor != 0 {
		query.Where("pagination > ?", cursor)
	}

	query.OrderExpr("pagination ASC")

	if err := query.Scan(ctx.Request().Context()); err != nil {
		return nil, strconv.FormatInt(0, 10), fmt.Errorf("selectNextPage: scan error: %w", err)
	}

	var newCursor string
	if len(products) > limit {
		newCursor = strconv.FormatInt(products[len(products)-2].Pagination, 10)
		products = products[:limit]
	} else {
		newCursor = ""
	}

	return products, newCursor, nil
}
