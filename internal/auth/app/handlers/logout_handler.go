package handlers

import (
	"go-demo/internal/auth/app/services"
	"go-demo/internal/shared"

	"github.com/labstack/echo"
)

const AFTER_LOGOUT_REDIRECT_PATH = "/auth/login"

type logoutHandler struct {
	authService services.AuthService
	flashStore  shared.FlashStore
}

func NewLogoutHandler(as services.AuthService, fs shared.FlashStore) logoutHandler {
	return logoutHandler{
		authService: as,
		flashStore:  fs,
	}
}

func (h logoutHandler) Handler(c echo.Context) error {
	h.authService.Logout(c.Request(), c.Response())
	shared.NewFlashMessage("You have been logged out successfully", "success").
		AddToSession(h.flashStore, c.Request(), c.Response())
	return c.Redirect(302, AFTER_LOGOUT_REDIRECT_PATH)
}
