package middlewares

import (
	"errors"
	"fmt"
	"go-demo/internal/app"
	"go-demo/internal/app/services/authservice"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type AuthMiddleware struct {
	SessionStore authservice.SessionStore
	UserRepo     authservice.AuthUserRepository
}

func (am AuthMiddleware) WithUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(app.AppContext)

		session, err := am.SessionStore.Get(c.Request(), authservice.SESSION_NAME)
		if err != nil {
			cc.Error(echo.NewHTTPError(500, fmt.Errorf("LoadUserMiddleware: Failed to get session: %w", err)))
			return next(cc)
		}

		if session.IsNew {
			return next(cc)
		}

		id, ok := session.Values[authservice.SESSION_USER_ID_FIELD].(uuid.UUID)
		if !ok {
			return next(cc)
		}

		user, err := am.UserRepo.SelectUserByID(id)
		if err != nil {
			if errors.Is(err, authservice.ErrUserNotFound) {
				return next(cc)
			}

			c.Error(echo.NewHTTPError(500, fmt.Errorf("LoadUserMiddleware: Failed to select user: %w", err)))

			return next(cc)
		}

		authenticatedUser := &app.AuthenticatedUser{
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
		cc := c.(app.AppContext)

		if cc.User == nil {
			return cc.Redirect(302, "auth/login")
		}

		return next(cc)
	}
}

func (am AuthMiddleware) WithAdminRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(app.AppContext)

		if cc.User == nil || !cc.User.IsAdmin {
			return cc.Redirect(302, "auth/login")
		}

		return next(cc)
	}
}
