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

	// The SQL statement we want to execute
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use `Exec()` for queries that do NOT return rows
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertID() method on the result to get the ID
	// of our newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Returned ID has the type int64; convert it to int type before returning.
	return int(id), nil
}

// Return a specific snippet based on id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Return 10 most recent snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
