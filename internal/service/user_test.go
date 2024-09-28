package service

import (
	"context"
	"database/sql"
	"github.com/atran25/synckor/internal/sqlc"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
	"testing"
)

func setupDatabase(tb testing.TB) (UserService, func(tb testing.TB)) {
	dbConnection, err := sql.Open("sqlite", ":memory:")
	require.Nil(tb, err, "Starting in memory sqlite database")
	_, err = migrate.Exec(dbConnection, "sqlite3", &migrate.FileMigrationSource{
		Dir: "../../db/migrations",
	}, migrate.Up)
	assert.NoError(tb, err, "Applying migrations")

	us := UserService{
		DB:  dbConnection,
		Qtx: sqlc.New(dbConnection),
	}

	return us, func(tb testing.TB) {
		dbConnection.Close()
	}
}

func TestGetDocumentSyncProgress(t *testing.T) {
	t.Run("document exists for user", func(t *testing.T) {
		us, teardown := setupDatabase(t)
		defer teardown(t)

		d := sqlc.DocumentInformation{
			Hash:       "testdocument",
			Username:   "testuser",
			Progress:   "testprogress",
			Device:     "testdevice",
			DeviceID:   "testdeviceid",
			Percentage: 0.5,
		}
		err := us.UpdateSyncProgress(context.Background(), float32(d.Percentage), d.Username, d.Hash, d.Progress, d.Device, d.DeviceID)
		require.NoError(t, err, "Adding initial sync progress")

		documentInformation, err := us.GetDocumentSyncProgress(context.Background(), d.Username, d.Hash)
		assert.NoError(t, err, "Getting document sync progress")
		assert.Equal(t, d.Hash, documentInformation.Hash)
		assert.Equal(t, d.Username, documentInformation.Username)
		assert.Equal(t, d.Progress, documentInformation.Progress)
		assert.Equal(t, d.Device, documentInformation.Device)
		assert.Equal(t, d.DeviceID, documentInformation.DeviceID)
		assert.Equal(t, d.Percentage, documentInformation.Percentage)
	})
	t.Run("document does not exist for user", func(t *testing.T) {
		us, teardown := setupDatabase(t)
		defer teardown(t)

		documentInformation, err := us.GetDocumentSyncProgress(context.Background(), "testuser", "testdocument")
		assert.Error(t, err, "Getting document sync progress")
		assert.Empty(t, documentInformation)
	})
}

func TestUpdateSyncProgress(t *testing.T) {
	us, teardown := setupDatabase(t)
	defer teardown(t)

	d := sqlc.DocumentInformation{
		Hash:       "testdocument",
		Username:   "testuser",
		Progress:   "testprogress",
		Device:     "testdevice",
		DeviceID:   "testdeviceid",
		Percentage: 0.5,
	}
	err := us.UpdateSyncProgress(context.Background(), float32(d.Percentage), d.Username, d.Hash, d.Progress, d.Device, d.DeviceID)
	require.NoError(t, err, "Adding initial sync progress")

	documentInformation, err := us.GetDocumentSyncProgress(context.Background(), d.Username, d.Hash)
	assert.NoError(t, err, "Getting document sync progress")
	assert.Equal(t, d.Hash, documentInformation.Hash)
	assert.Equal(t, d.Username, documentInformation.Username)
	assert.Equal(t, d.Progress, documentInformation.Progress)
	assert.Equal(t, d.Device, documentInformation.Device)
	assert.Equal(t, d.DeviceID, documentInformation.DeviceID)
	assert.Equal(t, d.Percentage, documentInformation.Percentage)

	d2 := sqlc.DocumentInformation{
		Hash:       "testdocument",
		Username:   "testuser",
		Progress:   "testprogress2",
		Device:     "testdevice2",
		DeviceID:   "testdeviceid2",
		Percentage: 0.75,
	}
	err = us.UpdateSyncProgress(context.Background(), float32(d2.Percentage), d2.Username, d2.Hash, d2.Progress, d2.Device, d2.DeviceID)
	require.NoError(t, err, "Updating sync progress")

	documentInformation, err = us.GetDocumentSyncProgress(context.Background(), d2.Username, d2.Hash)
	assert.NoError(t, err, "Getting document sync progress")
	assert.Equal(t, d2.Hash, documentInformation.Hash)
	assert.Equal(t, d2.Username, documentInformation.Username)
	assert.Equal(t, d2.Progress, documentInformation.Progress)
	assert.Equal(t, d2.Device, documentInformation.Device)
	assert.Equal(t, d2.DeviceID, documentInformation.DeviceID)
	assert.Equal(t, d2.Percentage, documentInformation.Percentage)
}

func TestNewUserService(t *testing.T) {
	dbConnection, err := sql.Open("sqlite", ":memory:")
	us := NewUserService(dbConnection)
	require.Nil(t, err, "Starting in memory sqlite database")
	assert.NotNil(t, us, "Creating user service")
	assert.Equal(t, dbConnection, us.DB, "Setting database connection")
	assert.Equal(t, sqlc.New(dbConnection), us.Qtx, "Setting query context")
}

func TestCreateUser(t *testing.T) {
	t.Run("user does not exist", func(t *testing.T) {
		us, teardown := setupDatabase(t)
		defer teardown(t)

		user := sqlc.UserAccount{
			Username:     "testuser",
			PasswordHash: "testpassword",
		}

		err := us.CreateUser(context.Background(), user.Username, user.PasswordHash)
		require.NoError(t, err, "Creating user")

		userAccount, err := us.Qtx.GetUser(context.Background(), user.Username)
		require.NoError(t, err, "Getting user")
		assert.Equal(t, user.Username, userAccount.Username)
		assert.Equal(t, user.PasswordHash, userAccount.PasswordHash)
	})
	t.Run("user already exists", func(t *testing.T) {
		us, teardown := setupDatabase(t)
		defer teardown(t)

		user := sqlc.UserAccount{
			Username:     "testuser",
			PasswordHash: "testpassword",
		}

		err := us.CreateUser(context.Background(), user.Username, user.PasswordHash)
		require.NoError(t, err, "Creating user")

		err = us.CreateUser(context.Background(), user.Username, user.PasswordHash)
		assert.Error(t, err, "Creating user")
		assert.Equal(t, ErrUserExists, err)
	})
}

func TestAuthenticateUser(t *testing.T) {
	t.Run("user exists", func(t *testing.T) {
		us, teardown := setupDatabase(t)
		defer teardown(t)

		user := sqlc.UserAccount{
			Username:     "testuser",
			PasswordHash: "testpassword",
		}

		err := us.CreateUser(context.Background(), user.Username, user.PasswordHash)
		require.NoError(t, err, "Creating user")

		status, err := us.AuthenticateUser(context.Background(), user.Username, user.PasswordHash)
		assert.NoError(t, err, "Authenticating user")
		assert.True(t, status)
	})
	t.Run("user does not exist", func(t *testing.T) {
		us, teardown := setupDatabase(t)
		defer teardown(t)

		user := sqlc.UserAccount{
			Username:     "testuser",
			PasswordHash: "testpassword",
		}

		status, err := us.AuthenticateUser(context.Background(), user.Username, user.PasswordHash)
		assert.Error(t, err, "Authenticating user")
		assert.False(t, status)
	})
	t.Run("password does not match", func(t *testing.T) {
		us, teardown := setupDatabase(t)
		defer teardown(t)

		user := sqlc.UserAccount{
			Username:     "testuser",
			PasswordHash: "testpassword",
		}

		err := us.CreateUser(context.Background(), user.Username, user.PasswordHash)
		require.NoError(t, err, "Creating user")

		status, err := us.AuthenticateUser(context.Background(), user.Username, "wrongpassword")
		assert.Error(t, err, "Authenticating user")
		assert.False(t, status)
	})
}
