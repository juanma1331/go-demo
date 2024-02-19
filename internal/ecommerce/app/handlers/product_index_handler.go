package handlers

import (
	"fmt"

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

	if cc.Flash != nil && len(*cc.Flash) > 0 {
		// This is not working
		firstFlash := (*cc.Flash)[0]
		notifyTrigger := shared.HtmxTrigger{
			Name: "notify",
			Value: map[string]string{
				"message": firstFlash.Message,
				"type":    firstFlash.Type,
			},
		}

		err := shared.SetHtmxTriggers(cc.Response().Writer, notifyTrigger)
		if err != nil {
			return fmt.Errorf("HandleAddToCart: error setting htmx triggers: %w", err)
		}
	}

	return cc.RenderComponent(ecommerce.IndexPage())
}
