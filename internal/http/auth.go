package http

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var (
	authProviders map[string]AuthProviderFactory
)

func init() {
	authProviders = make(map[string]AuthProviderFactory)
}

// AuthProviderFactory is a mechanism to return a complete auth system
// ready to go.  Because this registers a set of routes, the only
// thing it can really do is return a bad object, so we check an error
// here to see if something goes wrong.
type AuthProviderFactory func(*echo.Group) error

// TokenClaims is a type for setting the claims in the authentication
// tokens.  The only required claims that have to be set are the name,
// and networks the holder is an admin of.
type TokenClaims map[string]interface{}

// RegisterAuthProvider is used by other providers to be able to
// register them into the system and create a subtree on the /auth/
// route.
func RegisterAuthProvider(n string, ap AuthProviderFactory) {
	if _, here := authProviders[n]; here {
		// Already registered...
		return
	}
	authProviders[n] = ap
}

func (s *Server) initializeAuthProviders() {
	log.Println("Initializing authentication providers")
	ag := s.e.Group("/auth")
	for provider, initializer := range authProviders {
		log.Printf("Initializing %s", provider)
		tg := ag.Group("/" + provider)
		if err := initializer(tg); err != nil {
			log.Printf("Failed to initialize authentication provider %s: %v", provider, err)
		}
	}
}

// AuthCreateToken returns a token with the provided claims in it.  At
// the minimum its expected that a
func AuthCreateToken(c map[string]interface{}) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range c {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(viper.GetDuration("http.token.lifetime")).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("http.token.key")))
	return t, err
}

// requestAuthorized looks at the permissions in a given request and
// checks that the request either contains the requested permission,
// or a supervisory capability that overrides it.
func (s *Server) requestAuthorized(c echo.Context, perm, tgtOwner string) error {
	claims := getClaims(c)
	reqOwner := claims["user"].(string)

	pset := make(map[string]struct{})
	for _, p := range claims["permissions"].([]interface{}) {
		pset[strings.ToLower(p.(string))] = struct{}{}
	}
	log.Println(pset)

	net := c.Param("id")

	// Make checks for the individual nodes of the permissions
	// checks below.  This still winds up being O(N) and is much
	// easier to read than a clever loop.
	_, netAny := pset[net+":*"]
	_, netExact := pset[net+":"+perm]
	_, netSudo := pset[net+":sudo"]
	_, anyAny := pset["*:*"]
	_, anyExact := pset["*:"+perm]
	_, anySudo := pset["*:sudo"]

	hasPerm := netAny || netExact || anyAny || anyExact
	hasSudo := netAny || anyAny || netSudo || anySudo

	// Now we can perform the permissions check.  If the user has
	// any permission on any network, or any on the specific one
	// that satisfies the check right away.  Failing that we check
	// that they have the permission and optionally if they
	// require sudo.
	if hasPerm && (reqOwner == tgtOwner || hasSudo) {
		return nil
	}

	// Anything but the above and the request isn't authorized.
	return errors.New("requestor unqualified")
}

// Need to fish out the claims in a single function so that the
// peerFromContext can use the same logic and error handling.
func getClaims(c echo.Context) map[string]interface{} {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}
