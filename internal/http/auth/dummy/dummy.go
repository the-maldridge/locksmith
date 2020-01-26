package dummy

import (
	"strings"

	"github.com/labstack/echo"

	"github.com/the-maldridge/locksmith/internal/http"
)

func init() {
	http.RegisterAuthProvider("dummy", initialize)
}

type dummyAuth struct{}

func initialize(g *echo.Group) error {
	x := new(dummyAuth)

	g.GET("/poke", x.poke)
	g.GET("/auth", x.auth)
	return nil
}

func (da *dummyAuth) poke(c echo.Context) error {
	return c.String(200, "Don't poke me!\n")
}

func (da *dummyAuth) auth(c echo.Context) error {
	cl := http.TokenClaims{}

	cl["user"] = c.QueryParam("u")
	cl["permissions"] = strings.Split(c.QueryParam("p"), ",")

	t, err := http.AuthCreateToken(cl)
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.String(200, t+"\n")
}
