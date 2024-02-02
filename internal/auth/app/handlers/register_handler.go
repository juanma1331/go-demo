package handlers

import (
	"go-demo/internal/auth/app/services"
	"go-demo/internal/shared"
	"go-demo/views/auth"

	"github.com/labstack/echo"
)

const AFTER_REGISTER_REDIRECT_PATH = "/auth/login"

type registerHandler struct {
	authService services.AuthService
	flashStore  shared.FlashStore
}

func NewRegisterHandler(as services.AuthService, fs shared.FlashStore) registerHandler {
	return registerHandler{
		authService: as,
		flashStore:  fs,
	}
}

func (uh *registerHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)
	input := services.RegisterInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	output, err := uh.authService.Register(input)
	if err != nil {
		return err
	}

	if output.ValidationErrors != nil {
		vm := auth.RegisterPageViewModel{
			Errors: output.ValidationErrors,
		}

		return cc.RenderComponent(auth.RegisterPage(vm))
	}

	shared.NewFlashMessage("You have been registered successfully", "success").
		AddToSession(uh.flashStore, c.Request(), c.Response())
	return cc.Redirect(302, AFTER_REGISTER_REDIRECT_PATH)
}
