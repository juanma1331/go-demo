package shared

import (
	"github.com/labstack/echo"
)

type FlashMiddleware struct {
	FlashStore FlashStore
}

func (fm FlashMiddleware) WithFlashMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(AppContext)

		flashMessages, err := fm.FlashStore.LoadFlash(cc)
		if err != nil {
			cc.Error(echo.ErrInternalServerError)
			return next(cc)
		}

		cc.Flash = flashMessages

		return next(cc)
	}
}
