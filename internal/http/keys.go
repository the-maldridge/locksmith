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
	peer, err := peerFromContext(c)
	if err != nil {
		return err
	}

	if err := s.nm.AttemptNetworkRegistration(c.Param("id"), peer); err != nil {
		return c.JSON(http.StatusPreconditionFailed, err.Error())
	}

	return c.JSON(http.StatusOK, peer)
}

func (s *Server) approvePeer(c echo.Context) error {
	peer, err := peerFromContext(c)
	if err != nil {
		return err
	}

	if err := s.nm.ApprovePeer(c.Param("id"), peer.PubKey); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) disapprovePeer(c echo.Context) error {
	peer, err := peerFromContext(c)
	if err != nil {
		return err
	}

	if err := s.nm.DisapprovePeer(c.Param("id"), peer.PubKey); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) activatePeer(c echo.Context) error {
	peer, err := peerFromContext(c)
	if err != nil {
		return err
	}

	if err := s.nm.ActivatePeer(c.Param("id"), peer.PubKey); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deactivatePeer(c echo.Context) error {
	peer, err := peerFromContext(c)
	if err != nil {
		return err
	}

	if err := s.nm.DeactivatePeer(c.Param("id"), peer.PubKey); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func peerFromContext(c echo.Context) (models.Peer, error) {
	peer := models.Peer{}
	if err := c.Bind(&peer); err != nil {
		return models.Peer{}, err
	}

	return peer, nil
}
