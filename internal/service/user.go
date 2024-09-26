package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/atran25/synckor/internal/sqlc"
	"log/slog"
	"time"
)

type UserService struct {
	DB  *sql.DB
	Qtx *sqlc.Queries
}

var ErrUserExists = errors.New("user already exists")

type UserServiceInterface interface {
	CreateUser(ctx context.Context, username, password string) error
	AuthenticateUser(ctx context.Context, username, password string) (bool, error)
	UpdateSyncProgress(ctx context.Context, percentage float32, username, document, progress, device, deviceID string) error
	GetDocumentSyncProgress(ctx context.Context, username, documentHash string) (sqlc.DocumentInformation, error)
}

func (us *UserService) GetDocumentSyncProgress(ctx context.Context, username, documentHash string) (sqlc.DocumentInformation, error) {
	tx, err := us.DB.Begin()
	if err != nil {
		return sqlc.DocumentInformation{}, err
	}
	defer tx.Rollback()
	qtx := us.Qtx.WithTx(tx)

	documentInformation, err := qtx.GetDocument(ctx, sqlc.GetDocumentParams{
		Hash:     documentHash,
		Username: username,
	})
	if err != nil {
		return sqlc.DocumentInformation{}, err
	}

	err = tx.Commit()
	if err != nil {
		return sqlc.DocumentInformation{}, err
	}

	slog.Info("Retrieved document information", "document", documentInformation)
	return documentInformation, nil
}

func (us *UserService) UpdateSyncProgress(ctx context.Context, percentage float32, username, document, progress, device, deviceID string) error {
	tx, err := us.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := us.Qtx.WithTx(tx)

	_, err = qtx.GetDocument(ctx, sqlc.GetDocumentParams{
		Hash:     document,
		Username: username,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = qtx.CreateDocument(ctx, sqlc.CreateDocumentParams{
				Hash:       document,
				Progress:   progress,
				Percentage: float64(percentage),
				Device:     device,
				DeviceID:   deviceID,
				Timestamp:  time.Now(),
				Username:   username,
			})
			if err != nil {
				return err
			}
			err = tx.Commit()
			if err != nil {
				return err
			}
			slog.Info("Created document information", "document", document)
			return nil
		} else {
			return err
		}
	}

	documentInformation, err := qtx.UpdateDocument(ctx, sqlc.UpdateDocumentParams{
		Progress:   progress,
		Percentage: float64(percentage),
		Device:     device,
		DeviceID:   deviceID,
		Timestamp:  time.Now(),
		Hash:       document,
		Username:   username,
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	slog.Info("Updated document information", "document", documentInformation)
	return nil
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB:  db,
		Qtx: sqlc.New(db),
	}
}

func (us *UserService) CreateUser(ctx context.Context, username, password string) error {
	tx, err := us.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := us.Qtx.WithTx(tx)

	_, err = qtx.GetUser(ctx, username)
	if err == nil {
		return ErrUserExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     username,
		PasswordHash: password,
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) AuthenticateUser(ctx context.Context, username, password string) (bool, error) {
	tx, err := us.DB.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	qtx := us.Qtx.WithTx(tx)

	_, err = qtx.GetUserWithPassword(ctx, sqlc.GetUserWithPasswordParams{
		Username:     username,
		PasswordHash: password,
	})
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
