package shared

import (
	"github.com/labstack/echo"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
}
