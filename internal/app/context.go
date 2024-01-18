package app

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

type AppContext struct {
	echo.Context
	User  *AuthenticatedUser
	Flash *[]FlashMessage
}

type AuthenticatedUser struct {
	Email   string
	IsAdmin bool
}

type userKey string

var ContextUserKey userKey = "user"

func (c *AppContext) RenderComponent(component templ.Component) error {

	ctx := context.WithValue(c.Request().Context(), ContextUserKey, c.User)

	return component.Render(ctx, c.Response())
}
