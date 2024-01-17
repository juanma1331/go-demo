package handlers

import (
	"errors"
	"go-demo/internal/app"
	"go-demo/internal/app/services"
	"go-demo/view/authview"

	"github.com/labstack/echo"
)

const (
	AFTER_LOGIN_REDIRECT_PATH    = "/"
	AFTER_LOGOUT_REDIRECT_PATH   = "/auth/login"
	AFTER_REGISTER_REDIRECT_PATH = "/auth/login"
)

type authHandler struct {
	authService services.AuthService
	flashStore  app.FlashStore
}

func NewUserHandler(as services.AuthService, fs app.FlashStore) *authHandler {
	return &authHandler{
		authService: as,
		flashStore:  fs,
	}
}

func (uh *authHandler) HandleShowLogin(c echo.Context) error {
	cc := c.(app.AppContext)
	return cc.RenderComponent(authview.ShowLogin())
}

func (uh *authHandler) HandleLogin(c echo.Context) error {
	cc := c.(app.AppContext)

	input := services.LoginInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	if err := uh.authService.Login(cc.Request(), cc.Response(), input); err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return cc.RenderComponent(authview.ShowLogin())
		}

	}

	app.NewFlashMessage("You have been logged in successfully", "success").AddToSession(uh.flashStore, c.Request(), c.Response())
	return c.Redirect(302, AFTER_LOGIN_REDIRECT_PATH)
}

func (uh *authHandler) HandleShowRegister(c echo.Context) error {
	cc := c.(app.AppContext)
	return cc.RenderComponent(authview.ShowRegister())
}

func (uh *authHandler) HandleRegister(c echo.Context) error {
	cc := c.(app.AppContext)
	input := services.RegisterInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	output, err := uh.authService.Register(input)
	if err != nil {
		return err
	}

	if len(output.ValidationErrors) > 0 {
		// Here we render the register page again with the validation errors
		return cc.RenderComponent(authview.ShowRegister())
	}

	app.NewFlashMessage("You have been registered successfully", "success").AddToSession(uh.flashStore, c.Request(), c.Response())
	return cc.Redirect(302, AFTER_REGISTER_REDIRECT_PATH)
}

func (uh *authHandler) HandleLogout(c echo.Context) error {
	uh.authService.Logout(c.Request(), c.Response())
	app.NewFlashMessage("You have been logged out successfully", "success").AddToSession(uh.flashStore, c.Request(), c.Response())
	return c.Redirect(302, AFTER_LOGOUT_REDIRECT_PATH)
}
