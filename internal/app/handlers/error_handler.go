package handlers

import (
	"github.com/labstack/echo"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {

	// if he, ok := err.(*echo.HTTPError); ok {
	// 	code = he.Code
	// }

	c.Logger().Error(err)
}
