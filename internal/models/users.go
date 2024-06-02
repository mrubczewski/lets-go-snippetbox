package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	var id int

	statement := "INSERT INTO users (name, email, password, createdAt) VALUES ($1, $2, $3, $4) RETURNING id"
	err := m.DB.QueryRow(context.Background(), statement, name, email, password, time.Now()).Scan(&id)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
