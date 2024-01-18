package handlers

import (
	"errors"
	"go-demo/internal/app"
	"go-demo/internal/app/services/authservice"
	"go-demo/views/authview"

	"github.com/labstack/echo"
)

const (
	AFTER_LOGIN_REDIRECT_PATH    = "/"
	AFTER_LOGOUT_REDIRECT_PATH   = "/auth/login"
	AFTER_REGISTER_REDIRECT_PATH = "/auth/login"
)

type authHandler struct {
	authService authservice.AuthService
	flashStore  app.FlashStore
}

func NewAuthHandler(as authservice.AuthService, fs app.FlashStore) *authHandler {
	return &authHandler{
		authService: as,
		flashStore:  fs,
	}
}

func (uh *authHandler) HandleShowLogin(c echo.Context) error {
	cc := c.(app.AppContext)

	viewModel := authview.LoginPageViewModel{
		HasInvalidCredentials: false,
	}

	return cc.RenderComponent(authview.LoginPage(viewModel))
}

func (uh *authHandler) HandleLogin(c echo.Context) error {
	cc := c.(app.AppContext)

	input := authservice.LoginInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	if err := uh.authService.Login(cc.Request(), cc.Response(), input); err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			viewModel := authview.LoginPageViewModel{
				HasInvalidCredentials: true,
			}

			return cc.RenderComponent(authview.LoginPage(viewModel))
		}

	}

	app.NewFlashMessage("You have been logged in successfully", "success").
		AddToSession(uh.flashStore, c.Request(), c.Response())

	return c.Redirect(302, AFTER_LOGIN_REDIRECT_PATH)
}

func (uh *authHandler) HandleShowRegister(c echo.Context) error {
	cc := c.(app.AppContext)
	return cc.RenderComponent(authview.RegisterPage())
}

func (uh *authHandler) HandleRegister(c echo.Context) error {
	cc := c.(app.AppContext)
	input := authservice.RegisterInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	output, err := uh.authService.Register(input)
	if err != nil {
		return err
	}

	if output.ValidationErrors != nil {
		// Here we render the register page again with the validation errors
		return cc.RenderComponent(authview.RegisterPage())
	}

	app.NewFlashMessage("You have been registered successfully", "success").
		AddToSession(uh.flashStore, c.Request(), c.Response())
	return cc.Redirect(302, AFTER_REGISTER_REDIRECT_PATH)
}

func (uh *authHandler) HandleLogout(c echo.Context) error {
	uh.authService.Logout(c.Request(), c.Response())
	app.NewFlashMessage("You have been logged out successfully", "success").
		AddToSession(uh.flashStore, c.Request(), c.Response())
	return c.Redirect(302, AFTER_LOGOUT_REDIRECT_PATH)
}
