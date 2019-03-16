package http

import (
	"github.com/labstack/echo"
)

// Server encapsulates the components of the http server and its
// associated key management.
type Server struct {
	e *echo.Echo
}
