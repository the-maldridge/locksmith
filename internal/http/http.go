package http

import (
	"net/http"

	"github.com/labstack/echo"
)

// New initializes and returns an http.Server.
func New() (Server, error) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return Server{
		e,
	}, nil
}

// Serve is called to commence serving.  It will only return with an
// error if the serve call aborts.
func (s *Server) Serve() error {
	return s.e.Start(":1323")
}
