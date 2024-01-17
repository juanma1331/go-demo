package middlewares

import (
	"go-demo/internal/app"

	"github.com/labstack/echo"
)

type FlashMiddleware struct {
	FlashStore app.FlashStore
}

func (fm FlashMiddleware) LoadFlashMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(app.AppContext)

		flashMessages, err := fm.FlashStore.LoadFlash(cc)
		if err != nil {
			c.Error(echo.ErrInternalServerError)
			return next(c)
		}

		cc.Flash = flashMessages

		return next(cc)
	}
}
