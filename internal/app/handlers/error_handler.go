package handlers

import (
	"fmt"

	"github.com/labstack/echo"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {

	c.Logger().Error(err)
	fmt.Printf("CustomHTTPErrorHandler: %v\n", err)
}
