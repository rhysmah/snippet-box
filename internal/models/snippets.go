package models

import (
	"database/sql"
	"time"
)

// Define a snippet type to hold the data for an individual
// snippet. The fields correspond to the fields in the
// MySQL snippets table.
type Snippet struct {
	ID      int // Created automatically by DB
	Title   string
	Content string
	Created time.Time // Created automatically by DB
	Expires time.Time
}

// Define a SnippetModel type which warps an sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 0, nil
}

// Return a specific snippet based on id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Return 10 most recent snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
