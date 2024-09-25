package main

import (
	"context"
	"fmt"
	"github.com/atran25/synckor/internal/api"
	"github.com/atran25/synckor/internal/config"
	"github.com/atran25/synckor/internal/database"
	"github.com/atran25/synckor/internal/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type UserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
	Message    string `json:"message"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 402,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
		Message:        "This is it",
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	slog.Info("Config loaded: ", "Config", cfg)
	databaseConnection, err := database.GetConnection()
	if err != nil {
		panic(err)
	}
	db := sqlc.New(databaseConnection)
	slog.Info("Database connection established", "DB", db)
	server := api.NewServer(cfg, databaseConnection)
	r := chi.NewRouter()
	swagger, err := api.GetSwagger()
	swagger.Servers = nil
	if err != nil {
		panic(err)
	}
	r.Use(nethttpmiddleware.OapiRequestValidator(swagger))
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ip := request.RemoteAddr
			httpMethod := request.Method
			httpPath := request.URL.Path
			ctx := context.WithValue(request.Context(), "ip", ip)
			ctx = context.WithValue(ctx, "httpMethod", httpMethod)
			ctx = context.WithValue(ctx, "httpPath", httpPath)
			slog.Info("Request received", "IP", ip, "Method", httpMethod, "Path", httpPath)
			handler.ServeHTTP(writer, request.WithContext(ctx))
		})
	})

	// Serve the openapi doc on /doc
	workDir, _ := os.Getwd()
	fs := http.FileServer(http.Dir(filepath.Join(workDir, "doc")))
	r.Get("/doc/*", http.StripPrefix("/doc", fs).ServeHTTP)
	r.Get("/openapi.yaml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, filepath.Join(workDir, "api.yaml"))
	})

	h := api.HandlerFromMux(api.NewStrictHandler(server, nil), r)
	slog.Info(fmt.Sprintf("Starting server on port %d", cfg.Port))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h,
	}

	// Get all routes and their middlewares
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
	slog.Error("Server error: ", s.ListenAndServe())
}
