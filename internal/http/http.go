package http

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/nm"
)

// New initializes and returns an http.Server.
func New(netman nm.NetworkManager) (Server, error) {
	e := echo.New()

	s := Server{
		e:  e,
		nm: netman,
	}

	v1 := e.Group("/v1")
	v1.Use(middleware.JWT([]byte(viper.GetString("http.token.key"))))

	v1.GET("/networks/:id", s.getNet)
	v1.POST("/networks/:id/peers", s.registerPeer)
	v1.GET("/networks/:id/peers/config/:peer", s.fetchPeerConfiguration)
	v1.POST("/networks/:id/peers/deregister", s.deregisterPeer)
	v1.POST("/networks/:id/peers/approve", s.approvePeer)
	v1.POST("/networks/:id/peers/disapprove", s.disapprovePeer)
	v1.POST("/networks/:id/peers/activate", s.activatePeer)
	v1.POST("/networks/:id/peers/deactivate", s.deactivatePeer)

	s.initializeAuthProviders()

	return s, nil
}

// Serve is called to commence serving.  It will only return with an
// error if the serve call aborts.
func (s *Server) Serve() error {
	return s.e.Start(":1323")
}
