package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutesExist(t *testing.T) {
	// Doesn't need to populate DB and datamodel since the routes are not using the DB...
	testApp := Config{}

	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routes := []string{"/authenticate"}
	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

// routes -> chi routes
// route -> route to check if exist or not
func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	// We should expect the route to be found
	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}
