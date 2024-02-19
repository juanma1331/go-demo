package handlers

import (
	"github.com/juanma1331/go-demo/internal/auth/app/services"
	"github.com/juanma1331/go-demo/internal/shared"
	"github.com/juanma1331/go-demo/views/components"

	"github.com/labstack/echo"
)

type validateRegisterPassword struct {
	authService services.AuthService
}

func NewValidateRegisterPasswordHandler(as services.AuthService) validateRegisterPassword {
	return validateRegisterPassword{
		authService: as,
	}
}

func (h validateRegisterPassword) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)
	input := services.ValidateRegisterPasswordInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	output, err := h.authService.ValidateRegisterPassword(input)
	if err != nil {
		return err
	}

	if output.ValidationErrors != nil {
		return cc.RenderComponent(components.ValidationErrors(output.ValidationErrors, "Password"))
	}

	return cc.String(200, "")
}
