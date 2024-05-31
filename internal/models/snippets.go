package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expiresAt time.Time) (int, error) {
	_, err := m.DB.Exec(context.Background(), "INSERT INTO snippets (title, content, \"createdAt\", \"expiresAt\") VALUES ($1, $2, now(), $3)", title, content, expiresAt)
	if err != nil {
		log.Fatalf("Insert failed: %v\n", err)
	}
	fmt.Println("Insert successful!")
	return 0, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
