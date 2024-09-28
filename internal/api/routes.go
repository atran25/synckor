package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"path/filepath"
)

func SetupRoutes(r *chi.Mux) (*chi.Mux, error) {
	// Serve the openapi doc on /doc
	workDir, _ := os.Getwd()
	fs := http.FileServer(http.Dir(filepath.Join(workDir, "doc")))
	r.Get("/doc/*", http.StripPrefix("/doc", fs).ServeHTTP)
	r.Get("/openapi.yaml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, filepath.Join(workDir, "api.yaml"))
	})

	return r, nil
}
