package http

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

func (s *Server) fetchPeerConfiguration(c echo.Context) error {
	netQ := c.Param("id")
	peerQ := c.Param("peer")

	peerKey, err := url.QueryUnescape(peerQ)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	cfg, err := s.nm.GenerateConfigForPeer(netQ, peerKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, cfg)
}
