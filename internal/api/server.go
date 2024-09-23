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

func (s *Server) PutSyncsProgress(ctx context.Context, request PutSyncsProgressRequestObject) (PutSyncsProgressResponseObject, error) {
	username := request.Params.XAuthUser
	passwordHash := request.Params.XAuthKey
	document := request.Body.Document
	progress := request.Body.Progress
	percentage := request.Body.Percentage
	device := request.Body.Device
	deviceID := request.Body.DeviceId
	slog.Info("Received sync progress", "username", username, "password", passwordHash, "document", document, "progress", progress, "percentage", percentage, "device", device, "deviceID", deviceID, "IP", ctx.Value("ip"))

	return PutSyncsProgress200JSONResponse{}, nil
}

func (s *Server) GetUsersAuth(ctx context.Context, request GetUsersAuthRequestObject) (GetUsersAuthResponseObject, error) {
	username := request.Params.XAuthUser
	passwordHash := request.Params.XAuthKey
	slog.Info("Authenticating user", "username", username, "password", passwordHash, "IP", ctx.Value("ip"))

	user, err := s.DB.GetUser(ctx, username)
	if err != nil {
		slog.Error("Failed to get user", "error", err, "username", username, "IP", ctx.Value("ip"))
		message := "User does not exist"
		return GetUsersAuth401JSONResponse{
			Message:  &message,
			UserName: &username,
		}, nil
	}

	if user.Passwordhash != passwordHash {
		slog.Info("Authentication failed", "username", username, "IP", ctx.Value("ip"))
		message := "Authentication failed: Invalid password"
		return GetUsersAuth401JSONResponse{
			Message:  &message,
			UserName: &username,
		}, nil
	}

	message := "Authentication successful: User is authorized"
	return GetUsersAuth200JSONResponse{
		Message:  &message,
		UserName: &username,
	}, nil
}

func (s *Server) GetHealthcheck(ctx context.Context, request GetHealthcheckRequestObject) (GetHealthcheckResponseObject, error) {
	message := "Server is up and running"
	return GetHealthcheck200JSONResponse{
		Message: &message,
	}, nil
}

func (s *Server) PostUsersCreate(ctx context.Context, request PostUsersCreateRequestObject) (PostUsersCreateResponseObject, error) {
	username := request.Body.Username
	password := request.Body.Password
	registrationEnabled := s.Cfg.RegistrationEnabled

	// Return early if registration is disabled
	if !registrationEnabled {
		slog.Info("Registration is disabled", "username", *username, "IP", ctx.Value("ip"))
		message := "Registration is disabled"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	// Check if user already exists
	user, err := s.DB.GetUser(ctx, *username)
	if err == nil {
		slog.Info("User already exists", "User", user, "IP", ctx.Value("ip"))
		message := "User already exists"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	// Create user
	user, err = s.DB.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     *username,
		Passwordhash: *password,
		Isactive:     true,
		Isadmin:      false,
	})
	if err != nil {
		slog.Error("Failed to create user", "error", err, "username", *username, "IP", ctx.Value("ip"))
		message := "Failed to create user"
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
