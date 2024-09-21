package main

import (
	"fmt"
	"github.com/atran25/synckor/internal/api"
	"github.com/atran25/synckor/internal/config"
	"github.com/atran25/synckor/internal/database"
	"github.com/atran25/synckor/internal/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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
	server := api.NewServer(cfg, db)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	workDir, _ := os.Getwd()
	fs := http.FileServer(filepath.Join(workDir, "web"))
	r.Handle("/doc/", http.StripPrefix("/doc", fs))
	h := api.HandlerFromMux(api.NewStrictHandler(server, nil), r)
	slog.Info(fmt.Sprintf("Starting server on port %d", cfg.Port))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h,
	}
	// 👇 the walking function 🚶‍♂️
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
	slog.Error("Server error: ", s.ListenAndServe())

	//r.Use(render.SetContentType(render.ContentTypeJSON))
	//r.Get("/healthcheck", func(writer http.ResponseWriter, request *http.Request) {
	//	writer.Write([]byte("Up and running"))
	//})
	//r.Route("/users", func(r chi.Router) {
	//	r.Post("/auth", func(writer http.ResponseWriter, request *http.Request) {
	//		writer.Write([]byte("Authenticating user"))
	//	})
	//	r.Post("/create", func(writer http.ResponseWriter, request *http.Request) {
	//		var userPayload UserPayload
	//		err := json.NewDecoder(request.Body).Decode(&userPayload)
	//		slog.Info("", "userPayload", userPayload)
	//		if err != nil {
	//			slog.Error("decoding user payload: ", err)
	//			//writer.WriteHeader(http.StatusPaymentRequired)
	//			//writer.Write([]byte("Couldn't decode user payload"))
	//			render.Render(writer, request, ErrRender(err))
	//			return
	//		}
	//		ctx := context.Background()
	//		user, err := db.GetUser(ctx, userPayload.Username)
	//		if err == nil {
	//			slog.Error("user already exists", err)
	//			//writer.WriteHeader(http.StatusPaymentRequired)
	//			//writer.Write([]byte("User already exists"))
	//			render.Render(writer, request, ErrRender(errors.New("User already exists")))
	//			return
	//		}
	//
	//		user, err = db.CreateUser(ctx, sqlc.CreateUserParams{
	//			Username:     userPayload.Username,
	//			Passwordhash: userPayload.Password,
	//			Isactive:     true,
	//			Isadmin:      true,
	//		})
	//		if err != nil {
	//			slog.Error("creating user: ", err)
	//			writer.WriteHeader(http.StatusPaymentRequired)
	//			writer.Write([]byte("Couldn't create user"))
	//			return
	//		}
	//		slog.Info("User created: ", "User", user)
	//		writer.WriteHeader(http.StatusCreated)
	//		writer.Write([]byte("User created"))
	//	})
	//})
	//slog.Info("Starting server on port 8050")
	//http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
