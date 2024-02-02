package handlers

import (
	"go-demo/internal/shared"
	"go-demo/views/auth"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
)

type showRegisterHandler struct {
}

func NewShowRegisterHandler() showRegisterHandler {
	return showRegisterHandler{}
}

func (uh *showRegisterHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)

	csrfToken := csrf.Token(cc.Request())

	return cc.RenderComponent(auth.RegisterPage(auth.RegisterPageViewModel{
		CSRFToken: csrfToken,
	}))
}
