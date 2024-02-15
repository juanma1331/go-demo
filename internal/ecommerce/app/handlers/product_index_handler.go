package handlers

import (
	"github.com/juanma1331/go-demo/internal/shared"
	"github.com/juanma1331/go-demo/views/ecommerce"

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
