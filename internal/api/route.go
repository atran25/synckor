package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SetupRoutes() {
	r := chi.NewRouter()
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
