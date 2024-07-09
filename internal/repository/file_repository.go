package repository

import (
	"context"
	"database/sql"

	"github.com/miguel-martins/multicloud-storage-go/internal/models"
	"github.com/miguel-martins/multicloud-storage-go/internal/util"
)

type FileRepository struct {
	DB *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{DB: db}
}

func (ur *FileRepository) Save(user *models.User) error {
	_, err := ur.DB.ExecContext(context.Background(), "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (ur *FileRepository) Authenticate(username, password string) (bool, error) {
	var storedPassword string
	err := ur.DB.QueryRowContext(context.Background(), "SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, err
	}

	// Compare hashed passwords
	if util.ComparePasswords(storedPassword, password) != nil {
		return true, nil
	}
	return false, nil
}
