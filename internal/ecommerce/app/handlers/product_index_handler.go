package handlers

import (
	"go-demo/internal/shared"
	"go-demo/views/ecommerce"

	"github.com/labstack/echo"
)

type showProductIndexHandler struct{}

func NewShowProductIndexHandler() showProductIndexHandler {
	return showProductIndexHandler{}
}

func (h showProductIndexHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)
	return cc.RenderComponent(ecommerce.IndexPage())
}
