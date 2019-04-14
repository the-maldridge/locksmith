package http

import (
	"github.com/labstack/echo"

	"github.com/the-maldridge/locksmith/internal/nm"
)

// New initializes and returns an http.Server.
func New(netman nm.NetworkManager) (Server, error) {
	e := echo.New()

	s := Server{
		e:  e,
		nm: netman,
	}

	e.GET("/v1/network/:id", s.getNet)
	e.POST("/v1/network/:id/peers", s.registerPeer)
	e.POST("/v1/network/:id/peers/:peer/approve", s.approvePeer)
	e.POST("/v1/network/:id/peers/:peer/disapprove", s.disapprovePeer)
	e.POST("/v1/network/:id/peers/:peer/activate", s.activatePeer)
	e.POST("/v1/network/:id/peers/:peer/deactivate", s.deactivatePeer)

	return s, nil
}

// Serve is called to commence serving.  It will only return with an
// error if the serve call aborts.
func (s *Server) Serve() error {
	return s.e.Start(":1323")
}
