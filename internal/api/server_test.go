package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/atran25/synckor/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/testutil"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
	"testing"
)

func TestGetHealthcheck(t *testing.T) {
	dbConnection, err := sql.Open("sqlite", ":memory:")
	require.Nil(t, err, "Starting in memory sqlite database")
	defer dbConnection.Close()

	cfg := config.Config{}
	server := NewServer(cfg, dbConnection)
	r := chi.NewRouter()
	_ = HandlerFromMux(NewStrictHandler(server, nil), r)

	httpCall := testutil.NewRequest().Get("/healthcheck").GoWithHTTPHandler(t, r).Recorder
	assert.Equal(t, 200, httpCall.Code)

	var response Response
	err = json.NewDecoder(httpCall.Body).Decode(&response)
	require.NoError(t, err, "Decoding response")
	assert.Equal(t, "Server is up and running", *response.Message)
}

func TestPostUsersCreate(t *testing.T) {
	dbConnection, err := sql.Open("sqlite", ":memory:")
	require.Nil(t, err, "Starting in memory sqlite database")
	defer dbConnection.Close()
	_, err = migrate.Exec(dbConnection, "sqlite3", &migrate.FileMigrationSource{
		Dir: "../../db/migrations",
	}, migrate.Up)
	assert.NoError(t, err, "Applying migrations")

	cfg := config.Config{}
	server := NewServer(cfg, dbConnection)
	r := chi.NewRouter()
	_ = HandlerFromMux(NewStrictHandler(server, nil), r)

	t.Run("Registration is disabled", func(t *testing.T) {
		server.Cfg.RegistrationEnabled = false
		userName := "disabled"
		password := "disabled"
		request := PostUsersCreateJSONRequestBody{
			Username: &userName,
			Password: &password,
		}

		httpCall := testutil.NewRequest().Post("/users/create").WithJsonBody(request).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, 402, httpCall.Code)

		var response PostUsersCreate402JSONResponse
		err := json.NewDecoder(httpCall.Body).Decode(&response)
		require.NoError(t, err, "Decoding response")
		assert.Equal(t, "Registration is disabled", *response.Message)
	})

	t.Run("User already exists", func(t *testing.T) {
		server.Cfg.RegistrationEnabled = true
		userName := "test"
		password := "test"
		request := PostUsersCreateJSONRequestBody{
			Username: &userName,
			Password: &password,
		}

		err := server.UserService.CreateUser(context.Background(), userName, password)
		require.NoError(t, err, "Creating user")

		httpCall := testutil.NewRequest().Post("/users/create").WithJsonBody(request).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, 402, httpCall.Code)

		var response PostUsersCreate402JSONResponse
		err = json.NewDecoder(httpCall.Body).Decode(&response)
		require.NoError(t, err, "Decoding response")
		assert.Equal(t, "User already exists", *response.Message)
	})

	t.Run("User created successfully", func(t *testing.T) {
		server.Cfg.RegistrationEnabled = true
		userName := "new"
		password := "new"
		request := PostUsersCreateJSONRequestBody{
			Username: &userName,
			Password: &password,
		}

		httpCall := testutil.NewRequest().Post("/users/create").WithJsonBody(request).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, 201, httpCall.Code)

		var response PostUsersCreate201JSONResponse
		err := json.NewDecoder(httpCall.Body).Decode(&response)
		require.NoError(t, err, "Decoding response")
		assert.Equal(t, "User created successfully", *response.Message)
	})
}
