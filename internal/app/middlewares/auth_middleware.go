package middlewares

import (
	"errors"
	"fmt"
	"go-demo/internal/app"
	"go-demo/internal/app/services/authservice"

	"github.com/labstack/echo"
)

type AuthMiddleware struct {
	SessionStore     authservice.SessionStore
	UserRepo         authservice.AuthUserRepository
	AuthTokenRepo    authservice.AuthTokenRepository
	AuthTokenManager authservice.AuthTokenManager
}

func (am AuthMiddleware) LoadUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		token, err := am.AuthTokenManager.ExtractTokenFromSession(session)
		if err != nil {
			return next(cc)
		}

		authToken, err := am.AuthTokenRepo.SelectToken(token.String())
		if err != nil {
			if errors.Is(err, authservice.ErrTokenNotFound) {
				return next(cc)
			}

			c.Error(echo.NewHTTPError(500, fmt.Errorf("LoadUserMiddleware: Failed to select token: %w", err)))
			return next(cc)
		}

		user, err := am.UserRepo.SelectUserByID(authToken.UserID.String())
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
		}

		cc.User = authenticatedUser

		return next(cc)
	}
}

func (am AuthMiddleware) RequireLoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(app.AppContext)

		if cc.User == nil {
			return c.Redirect(302, "auth/login")
		}

		return next(cc)
	}
}

func (am AuthMiddleware) RequireAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(app.AppContext)

		if cc.User == nil || !cc.User.IsAdmin {
			return c.Redirect(302, "auth/login")
		}

		return next(cc)
	}
}
