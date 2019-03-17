package http

import (
	"log"
	"net/http"

	"github.com/labstack/echo"

	"github.com/the-maldridge/locksmith/internal/models"
)

func (s *Server) registerClient(c echo.Context) error {
	// Check if the network requested is actually known.
	if _, err := s.nm.GetNet(c.Param("id")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// If the network is known, then figure out the request.
	client := new(models.Client)
	if err := c.Bind(client); err != nil {
		return err
	}

	log.Println("Network:", c.Param("id"), client)

	if err := s.nm.AttemptNetworkRegistration(c.Param("id"), *client); err != nil {
		return c.JSON(http.StatusPreconditionFailed, err.Error())
	}

	return c.JSON(http.StatusOK, client)
}
