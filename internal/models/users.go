package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	hashedPassword := string(hash)

	statement := "INSERT INTO users (name, email, hashed_password, \"createdAt\") VALUES ($1, $2, $3, $4) RETURNING id"
	err = m.DB.QueryRow(context.Background(), statement, name, email, hashedPassword, time.Now()).Scan(&id)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

func (m *UserModel) EmailTaken(email string) (bool, error) {
	statement := "SELECT id FROM users WHERE email = $1"
	row := m.DB.QueryRow(context.Background(), statement, email)

	var existingUserId int
	err := row.Scan(&existingUserId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	if existingUserId != 0 {
		return true, nil
	}
	return false, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"
	err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
