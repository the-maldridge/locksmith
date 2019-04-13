package http

import (
	"net/http"

	"github.com/labstack/echo"
)

func (s *Server) getNet(c echo.Context) error {
	// Check if the network requested is actually known.
	net, err := s.nm.GetNet(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, net)
}
