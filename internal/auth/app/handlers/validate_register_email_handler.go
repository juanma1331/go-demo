package handlers

import (
	"github.com/juanma1331/go-demo/internal/auth/app/services"
	"github.com/juanma1331/go-demo/internal/shared"
	"github.com/juanma1331/go-demo/views/components"

	"github.com/labstack/echo"
)

type validateRegisterEmailHandler struct {
	authService services.AuthService
}

func NewValidateRegisterEmailHandler(as services.AuthService) validateRegisterEmailHandler {
	return validateRegisterEmailHandler{
		authService: as,
	}
}

func (h validateRegisterEmailHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)
	input := services.ValidateRegisterEmailInput{}
	if err := cc.Bind(&input); err != nil {
		return err
	}

	output, err := h.authService.ValidateRegisterEmail(input)
	if err != nil {
		return err
	}

	if output.ValidationErrors != nil {
		errors := (*output.ValidationErrors)["Email"]

		return cc.RenderComponent(components.ValidationErrors(errors))
	}

	return cc.String(200, "")
}
