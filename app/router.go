package app

import (
	"net/http"

	"github.com/rlongo/ambrosia/api"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

// NewRouter exports a new router class and used Dependencu Injection to introduce
// any externally required items
func NewRouter(storage api.StorageServiceRecipes, middleware *negroni.Negroni) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	apiv1Router := router.PathPrefix("/api/v1/").Subrouter()
	for _, route := range ambrosiaRoutes {
		var handler http.Handler
		handler = route.HandlerFunc(storage)
		handler = middleware.With(
			negroni.Wrap(handler),
		)

		apiv1Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
