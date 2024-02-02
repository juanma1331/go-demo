package middlewares

import (
	"errors"
	"fmt"
	"go-demo/internal/auth/app/services"
	"go-demo/internal/shared"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type AuthMiddleware struct {
	SessionStore services.SessionStore
	UserRepo     services.AuthUserRepository
}

func (am AuthMiddleware) WithUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(shared.AppContext)

		session, err := am.SessionStore.Get(c.Request(), services.SESSION_NAME)
		if err != nil {
			cc.Error(echo.NewHTTPError(500, fmt.Errorf("WithUserMiddleware: Failed to get session: %w", err)))
			return next(cc)
		}

		if session.IsNew {
			return next(cc)
		}

		id, ok := session.Values[services.SESSION_USER_ID_FIELD].(uuid.UUID)
		if !ok {
			return next(cc)
		}

		user, err := am.UserRepo.SelectUserByID(id)
		if err != nil {
			if errors.Is(err, services.ErrUserNotFound) {
				return next(cc)
			}

			c.Error(echo.NewHTTPError(500, fmt.Errorf("WithUserMiddleware: Failed to select user: %w", err)))

			return next(cc)
		}

		authenticatedUser := &shared.AuthenticatedUser{
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
			ID:      user.ID,
		}

		cc.User = authenticatedUser

		return next(cc)
	}
}

func (am AuthMiddleware) WithAuthenticationRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(shared.AppContext)

		if cc.User == nil {
			return cc.Redirect(302, "auth/login")
		}

		return next(cc)
	}
}

func (am AuthMiddleware) WithAdminRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(shared.AppContext)

		if cc.User == nil || !cc.User.IsAdmin {
			return cc.Redirect(302, "auth/login")
		}

		return next(cc)
	}
}
