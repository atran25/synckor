package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/atran25/synckor/internal/config"
	"github.com/atran25/synckor/internal/service"
	"log/slog"
)

type Server struct {
	Cfg         config.Config
	UserService service.UserServiceInterface
}

func (s *Server) GetSyncsProgressDocumentHash(ctx context.Context, request GetSyncsProgressDocumentHashRequestObject) (GetSyncsProgressDocumentHashResponseObject, error) {
	documentHash := request.DocumentHash
	username := request.Params.XAuthUser
	passwordHash := request.Params.XAuthKey

	status, err := s.UserService.AuthenticateUser(ctx, username, passwordHash)
	if err != nil {
		slog.Error("Failed to authenticate user", "error", err, "username", username, "IP", ctx.Value("ip"))
		message := "Failed to authenticate user"
		return GetSyncsProgressDocumentHash401JSONResponse{
			Message: &message,
		}, nil
	}
	if !status {
		slog.Info("Authentication failed", "username", username, "IP", ctx.Value("ip"))
		message := "Authentication failed: User is not authorized"
		return GetSyncsProgressDocumentHash401JSONResponse{
			Message: &message,
		}, nil
	}

	documentInformation, err := s.UserService.GetDocumentSyncProgress(ctx, username, documentHash)
	if err != nil {
		slog.Error("Failed to get document sync progress", "error", err, "username", username, "documentHash", documentHash, "IP", ctx.Value("ip"))
		message := "Failed to get document sync progress"
		return GetSyncsProgressDocumentHash401JSONResponse{
			Message: &message,
		}, nil
	}

	// Workaround for converting float64 to *float32
	percentageConvert := float32(documentInformation.Percentage)
	return GetSyncsProgressDocumentHash200JSONResponse{
		Document:   &documentInformation.Hash,
		Progress:   &documentInformation.Progress,
		Percentage: &percentageConvert,
		Device:     &documentInformation.Device,
		DeviceId:   &documentInformation.DeviceID,
	}, nil
}

func (s *Server) PutSyncsProgress(ctx context.Context, request PutSyncsProgressRequestObject) (PutSyncsProgressResponseObject, error) {
	username := request.Params.XAuthUser
	passwordHash := request.Params.XAuthKey
	document := request.Body.Document
	progress := request.Body.Progress
	percentage := request.Body.Percentage
	device := request.Body.Device
	deviceID := request.Body.DeviceId

	status, err := s.UserService.AuthenticateUser(ctx, username, passwordHash)
	if err != nil {
		slog.Error("Failed to authenticate user", "error", err, "username", username, "IP", ctx.Value("ip"))
		message := "Failed to authenticate user"
		return PutSyncsProgress401JSONResponse{
			Message: &message,
		}, nil
	}
	if !status {
		slog.Info("Authentication failed", "username", username, "IP", ctx.Value("ip"))
		message := "Authentication failed: User is not authorized"
		return PutSyncsProgress401JSONResponse{
			Message: &message,
		}, nil
	}

	if document == nil || progress == nil || percentage == nil || device == nil || deviceID == nil {
		slog.Info("Invalid request", "username", username, "password", passwordHash, "document", document, "progress", progress, "percentage", percentage, "device", device, "deviceID", deviceID, "IP", ctx.Value("ip"))
		message := "Invalid request"
		return PutSyncsProgress401JSONResponse{
			Message: &message,
		}, nil
	}

	slog.Info("Received sync progress", "username", username, "password", passwordHash, "document", *document, "progress", *progress, "percentage", *percentage, "device", *device, "deviceID", *deviceID, "IP", ctx.Value("ip"))

	err = s.UserService.UpdateSyncProgress(ctx, *percentage, username, *document, *progress, *device, *deviceID)
	if err != nil {
		slog.Error("Failed to update sync progress", "error", err, "username", username, "document", *document, "progress", *progress, "percentage", *percentage, "device", *device, "deviceID", *deviceID, "IP", ctx.Value("ip"))
		message := "Failed to update sync progress"
		return PutSyncsProgress401JSONResponse{
			Message: &message,
		}, nil
	}

	message := "Sync progress updated successfully"
	return PutSyncsProgress200JSONResponse{
		Message: &message,
	}, nil
}

func (s *Server) GetUsersAuth(ctx context.Context, request GetUsersAuthRequestObject) (GetUsersAuthResponseObject, error) {
	username := request.Params.XAuthUser
	passwordHash := request.Params.XAuthKey
	slog.Info("Authenticating user", "username", username, "password", passwordHash, "IP", ctx.Value("ip"))

	status, err := s.UserService.AuthenticateUser(ctx, username, passwordHash)
	if err != nil {
		slog.Error("Failed to authenticate user", "error", err, "username", username, "IP", ctx.Value("ip"))
		message := "Failed to authenticate user"
		return GetUsersAuth401JSONResponse{
			Message:  &message,
			UserName: &username,
		}, nil
	}

	if !status {
		slog.Info("Authentication failed", "username", username, "IP", ctx.Value("ip"))
		message := "Authentication failed: User is not authorized"
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

	// Return early if username or password is nil
	if username == nil || password == nil {
		slog.Info("Invalid request", "username", username, "password", password, "IP", ctx.Value("ip"))
		message := "Invalid request"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	// Return early if registration is disabled
	if !registrationEnabled {
		slog.Info("Registration is disabled", "username", *username, "IP", ctx.Value("ip"))
		message := "Registration is disabled"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	err := s.UserService.CreateUser(ctx, *username, *password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			slog.Info("User already exists", "username", *username, "IP", ctx.Value("ip"))
			message := "User already exists"
			return PostUsersCreate402JSONResponse{
				Message: &message,
			}, nil
		} else {
			slog.Error("Failed to create user", "error", err, "username", *username, "IP", ctx.Value("ip"))
			message := "Failed to create user"
			return PostUsersCreate402JSONResponse{
				Message: &message,
			}, nil
		}
	}

	message := "User created successfully"
	return PostUsersCreate201JSONResponse{
		Message: &message,
	}, nil
}

func NewServer(cfg config.Config, DB *sql.DB) *Server {
	return &Server{
		Cfg:         cfg,
		UserService: service.NewUserService(DB),
	}
}
