package handlers

import (
	"errors"
	"fmt"
	"go-demo/internal/auth/app/services"
	"go-demo/internal/shared"
	"go-demo/views/auth"

	"github.com/labstack/echo"
)

const AFTER_LOGIN_REDIRECT_PATH = "/"

type loginHandler struct {
	authService services.AuthService
	flashStore  shared.FlashStore
}

func NewLoginHandler(as services.AuthService, fs shared.FlashStore) loginHandler {
	return loginHandler{
		authService: as,
		flashStore:  fs,
	}
}

func (uh *loginHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)

	input := services.LoginInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	if err := uh.authService.Login(cc.Request(), cc.Response(), input); err != nil {
		fmt.Println(err)
		if errors.Is(err, services.ErrInvalidCredentials) {
			viewModel := auth.LoginPageViewModel{
				HasInvalidCredentials: true,
			}

			return cc.RenderComponent(auth.LoginPage(viewModel))
		}

		return err
	}

	shared.NewFlashMessage("You have been logged in successfully", "success").
		AddToSession(uh.flashStore, c.Request(), c.Response())

	return c.Redirect(302, AFTER_LOGIN_REDIRECT_PATH)
}
