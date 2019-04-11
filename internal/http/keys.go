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
	netID, pubkey, err := s.parseKeyFromContext(c)
	if err != nil {
		return err
	}

	if err := s.nm.ApprovePeer(netID, pubkey); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) activatePeer(c echo.Context) error {
	netID, pubkey, err := s.parseKeyFromContext(c)
	if err != nil {
		return err
	}

	if err := s.nm.ActivatePeer(netID, pubkey); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) parseKeyFromContext(c echo.Context) (string, string, error) {
	// Check if the network requested is actually known.
	if _, err := s.nm.GetNet(c.Param("id")); err != nil {
		return "", "", c.String(http.StatusBadRequest, err.Error())
	}

	var b struct {
		PubKey string
	}
	if err := c.Bind(&b); err != nil {
		return "", "", err
	}

	return c.Param("id"), b.PubKey, nil
}
