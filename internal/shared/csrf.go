package shared

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo"
)

type CSRFMiddleware struct{}

const (
	CSRFTokenKey = "csrf_token"
)

func (CSRFMiddleware) WithCSRFMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	CSRF := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.Secure(false), // True in production
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.FieldName(CSRFTokenKey),
	)

	return func(c echo.Context) error {
		handler := CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.SetRequest(r)
			next(c)
		}))

		handler.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
