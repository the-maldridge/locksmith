package http

import (
	"github.com/labstack/echo"

	"github.com/the-maldridge/locksmith/internal/nm"
)

// Server encapsulates the components of the http server and its
// associated key management.
type Server struct {
	e  *echo.Echo
	nm nm.NetworkManager
}
