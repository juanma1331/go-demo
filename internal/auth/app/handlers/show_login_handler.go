package handlers

import (
	"github.com/juanma1331/go-demo/internal/shared"
	"github.com/juanma1331/go-demo/views/auth"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
)

type showLoginHandler struct {
}

func NewShowLoginHandler() showLoginHandler {
	return showLoginHandler{}
}

func (h showLoginHandler) Handler(c echo.Context) error {
	cc := c.(shared.AppContext)

	csrfToken := csrf.Token(cc.Request())

	viewModel := auth.LoginPageViewModel{
		HasInvalidCredentials: false,
		CSRFToken:             csrfToken,
	}

	return cc.RenderComponent(auth.LoginPage(viewModel))
}
