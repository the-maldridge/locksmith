package http

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/the-maldridge/locksmith/internal/models"
)

func (s *Server) registerPeer(c echo.Context) error {
	// Check if the network requested is actually known.
	if _, err := s.nm.GetNet(c.Param("id")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// If the network is known, then figure out the request.
	client := new(models.Peer)
	if err := c.Bind(client); err != nil {
		return err
	}

	if err := s.nm.AttemptNetworkRegistration(c.Param("id"), *client); err != nil {
		return c.JSON(http.StatusPreconditionFailed, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}

func (s *Server) approvePeer(c echo.Context) error {
	if err := s.nm.ApprovePeer(c.Param("id"), c.Param("peer")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) disapprovePeer(c echo.Context) error {
	if err := s.nm.DisapprovePeer(c.Param("id"), c.Param("peer")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) activatePeer(c echo.Context) error {
	if err := s.nm.ActivatePeer(c.Param("id"), c.Param("peer")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deactivatePeer(c echo.Context) error {
	if err := s.nm.DeactivatePeer(c.Param("id"), c.Param("peer")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
