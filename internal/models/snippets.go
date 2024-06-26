package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
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

func (m *SnippetModel) Insert(title string, content string, expireAfterDays int) (int, error) {
	var id int

	statement := "INSERT INTO snippets (title, content, \"createdAt\", \"expiresAt\") VALUES ($1, $2, $3, $4) RETURNING id"
	currentDateTime := time.Now()
	err := m.DB.QueryRow(context.Background(), statement, title, content, currentDateTime, currentDateTime.AddDate(0, 0, expireAfterDays)).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	statement := "SELECT id, title, content, \"createdAt\", \"expiresAt\" FROM snippets WHERE \"expiresAt\" > NOW() AND id = $1"
	row := m.DB.QueryRow(context.Background(), statement, id)
	var s Snippet
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)

	newContent := strings.Replace(s.Content, "\n", "<br>", -1)
	log.Printf(newContent)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	statement := "SELECT id, title, content, \"createdAt\", \"expiresAt\" FROM snippets WHERE \"expiresAt\" > NOW() ORDER BY \"createdAt\" DESC LIMIT 10"
	rows, err := m.DB.Query(context.Background(), statement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
