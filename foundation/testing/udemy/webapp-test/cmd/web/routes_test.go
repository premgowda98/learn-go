package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/static/*", "GET"},
	}

	app := application{}
	mux := app.routes()

	chiRoutes := mux.(chi.Routes)

	for _, route := range registered{
		if !routeExists(route.route, route.method, chiRoutes){
			t.Error("routes check failed")
		}
	}

	
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found:=false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(testRoute, route) && strings.EqualFold(method, testMethod){
			found = true
		}

		return nil
	})

	return found
}
