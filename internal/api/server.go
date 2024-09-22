package api

import (
	"context"
	"github.com/atran25/synckor/internal/config"
	"github.com/atran25/synckor/internal/sqlc"
	"log/slog"
)

type Server struct {
	Cfg config.Config
	DB  *sqlc.Queries
}

func (s *Server) GetHealthcheck(ctx context.Context, request GetHealthcheckRequestObject) (GetHealthcheckResponseObject, error) {
	message := "Server is up and running"
	return GetHealthcheck200JSONResponse{
		Message: &message,
	}, nil
}

func (s *Server) PostUsersCreate(ctx context.Context, request PostUsersCreateRequestObject) (PostUsersCreateResponseObject, error) {
	username := request.Body.Username
	//password := request.Body.Password
	registrationEnabled := s.Cfg.RegistrationEnabled

	// Return early if registration is disabled
	if !registrationEnabled {
		slog.Info("Registration is disabled", "username", *username, "IP", ctx.Value("ip"))
		message := "Registration is disabled"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	user, err := s.DB.GetUser(ctx, *username)
	if err == nil {
		slog.Info("User already exists", "User", user)
		message := "User already exists"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	message := "User created successfully"
	return PostUsersCreate201JSONResponse{
		Message: &message,
	}, nil
}

func NewServer(cfg config.Config, DB *sqlc.Queries) *Server {
	return &Server{
		Cfg: cfg,
		DB:  DB,
	}
}
