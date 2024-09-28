package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/atran25/synckor/internal/config"
	"github.com/atran25/synckor/internal/service"
	"github.com/go-chi/chi/v5"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"log/slog"
	"net/http"
)

type IP struct{}
type HttpMethod struct{}
type HttpPath struct{}

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
		slog.Error("Failed to authenticate user", "error", err, "username", username, "IP", ctx.Value(IP{}))
		message := "Failed to authenticate user"
		return GetSyncsProgressDocumentHash401JSONResponse{
			Message: &message,
		}, nil
	}
	if !status {
		slog.Info("Authentication failed", "username", username, "IP", ctx.Value(IP{}))
		message := "Authentication failed: User is not authorized"
		return GetSyncsProgressDocumentHash401JSONResponse{
			Message: &message,
		}, nil
	}

	documentInformation, err := s.UserService.GetDocumentSyncProgress(ctx, username, documentHash)
	if err != nil {
		slog.Error("Failed to get document sync progress", "error", err, "username", username, "documentHash", documentHash, "IP", ctx.Value(IP{}))
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
		slog.Error("Failed to authenticate user", "error", err, "username", username, "IP", ctx.Value(IP{}))
		message := "Failed to authenticate user"
		return PutSyncsProgress401JSONResponse{
			Message: &message,
		}, nil
	}
	if !status {
		slog.Info("Authentication failed", "username", username, "IP", ctx.Value(IP{}))
		message := "Authentication failed: User is not authorized"
		return PutSyncsProgress401JSONResponse{
			Message: &message,
		}, nil
	}

	if document == nil || progress == nil || percentage == nil || device == nil || deviceID == nil {
		slog.Info("Invalid request", "username", username, "password", passwordHash, "document", document, "progress", progress, "percentage", percentage, "device", device, "deviceID", deviceID, "IP", ctx.Value(IP{}))
		message := "Invalid request"
		return PutSyncsProgress401JSONResponse{
			Message: &message,
		}, nil
	}

	slog.Info("Received sync progress", "username", username, "password", passwordHash, "document", *document, "progress", *progress, "percentage", *percentage, "device", *device, "deviceID", *deviceID, "IP", ctx.Value(IP{}))

	err = s.UserService.UpdateSyncProgress(ctx, *percentage, username, *document, *progress, *device, *deviceID)
	if err != nil {
		slog.Error("Failed to update sync progress", "error", err, "username", username, "document", *document, "progress", *progress, "percentage", *percentage, "device", *device, "deviceID", *deviceID, "IP", ctx.Value(IP{}))
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
	slog.Info("Authenticating user", "username", username, "password", passwordHash, "IP", ctx.Value(IP{}))

	status, err := s.UserService.AuthenticateUser(ctx, username, passwordHash)
	if err != nil {
		slog.Error("Failed to authenticate user", "error", err, "username", username, "IP", ctx.Value(IP{}))
		message := "Failed to authenticate user"
		return GetUsersAuth401JSONResponse{
			Message:  &message,
			UserName: &username,
		}, nil
	}

	if !status {
		slog.Info("Authentication failed", "username", username, "IP", ctx.Value(IP{}))
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
		slog.Info("Invalid request", "username", username, "password", password, "IP", ctx.Value(IP{}))
		message := "Invalid request"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	// Return early if registration is disabled
	if !registrationEnabled {
		slog.Info("Registration is disabled", "username", *username, "IP", ctx.Value(IP{}))
		message := "Registration is disabled"
		return PostUsersCreate402JSONResponse{
			Message: &message,
		}, nil
	}

	err := s.UserService.CreateUser(ctx, *username, *password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			slog.Info("User already exists", "username", *username, "IP", ctx.Value(IP{}))
			message := "User already exists"
			return PostUsersCreate402JSONResponse{
				Message: &message,
			}, nil
		} else {
			slog.Error("Failed to create user", "error", err, "username", *username, "IP", ctx.Value(IP{}))
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

func NewServer(cfg config.Config, DB *sql.DB) (http.Handler, error) {
	server := &Server{
		Cfg:         cfg,
		UserService: service.NewUserService(DB),
	}
	r := chi.NewRouter()

	// Add request validator middleware
	swagger, err := GetSwagger()
	if err != nil {
		return nil, err
	}
	swagger.Servers = nil
	r.Use(nethttpmiddleware.OapiRequestValidator(swagger))
	r.Use(RequestLogger)

	r, err = SetupRoutes(r)
	if err != nil {
		return nil, err
	}

	h := HandlerFromMux(NewStrictHandler(server, nil), r)

	err = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		slog.Info("Route found", "method", method, "route", route, "middlewares", len(middlewares))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return h, nil
}

func RequestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ip := request.RemoteAddr
		httpMethod := request.Method
		httpPath := request.URL.Path
		ctx := context.WithValue(request.Context(), IP{}, ip)
		ctx = context.WithValue(ctx, HttpMethod{}, httpMethod)
		ctx = context.WithValue(ctx, HttpPath{}, httpPath)
		slog.Info("Request received", "IP", ip, "Method", httpMethod, "Path", httpPath)
		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}
