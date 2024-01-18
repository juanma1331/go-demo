package middlewares

import (
	"go-demo/internal/app"

	"github.com/labstack/echo"
)

type AppContextMiddleware struct{}

func (ap AppContextMiddleware) WithAppContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := app.AppContext{
			Context: c,
		}

		return next(cc)
	}
}
