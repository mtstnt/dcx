package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	PREFIX  = "/api"
	VERSION = "v1"
)

func Routes() chi.Router {
	router := chi.NewRouter()

	prefix := fmt.Sprintf("%s/%s", PREFIX, VERSION)
	router.Route(prefix, func(r chi.Router) {
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	})

	return router
}
