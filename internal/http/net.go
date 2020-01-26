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

	if err = s.requestAuthorized(c, "sudo", ""); err != nil {
		return c.JSON(http.StatusPreconditionFailed, err.Error())
	}

	return c.JSON(http.StatusOK, net)
}
