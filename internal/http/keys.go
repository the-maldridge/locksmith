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

func (s *Server) activateClient(c echo.Context) error {
	// Check if the network requested is actually known.
	if _, err := s.nm.GetNet(c.Param("id")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var b struct {
		PubKey string
	}
	if err := c.Bind(&b); err != nil {
		return err
	}

	// TODO: Dirty rotten hack to stick the equal sign back on the
	// end of the key after the JSON parser eats it.  Replace this
	// with something more intelligent later.
	b.PubKey += "="

	log.Println(b)

	if err := s.nm.ActivatePeer(c.Param("id"), b.PubKey); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
