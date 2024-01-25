package app

import (
	"context"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type AppContext struct {
	echo.Context
	User  *AuthenticatedUser
	Flash *[]FlashMessage
}

type AuthenticatedUser struct {
	ID      uuid.UUID
	Email   string
	IsAdmin bool
}

type contextKey string

const (
	ContextUserKey  contextKey = "user"
	ContextFlashKey contextKey = "flash"
)

func (c *AppContext) RenderComponent(component templ.Component) error {

	ctx := context.WithValue(c.Request().Context(), ContextUserKey, c.User)
	ctx = context.WithValue(ctx, ContextFlashKey, c.Flash)

	return component.Render(ctx, c.Response())
}
