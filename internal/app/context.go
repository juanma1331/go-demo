package app

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

type AppContext struct {
	echo.Context
	User  *AuthenticatedUser
	Flash []FlashMessage
}

type AuthenticatedUser struct {
	Email   string
	IsAdmin bool
}

func (c *AppContext) RenderComponent(component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

type AppContextMiddleware struct{}

func (ap AppContextMiddleware) CreateAppContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := AppContext{
			Context: c,
		}

		return next(cc)
	}
}
