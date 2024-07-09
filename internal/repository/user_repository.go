package repository

import (
	"context"
	"database/sql"

	"github.com/miguel-martins/multicloud-storage-go/internal/models"
	"github.com/miguel-martins/multicloud-storage-go/internal/util"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) Save(user *models.User) error {
	_, err := ur.DB.ExecContext(context.Background(), "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) Authenticate(username, password string) (bool, error) {
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
