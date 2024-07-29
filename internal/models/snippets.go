package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a snippet type to hold the data for an individual snippet.
// The fields correspond to the fields in the MySQL snippets table.
type Snippet struct {
	ID      int // Created automatically by DB
	Title   string
	Content string
	Created time.Time // Created automatically by DB
	Expires time.Time
}

// Define a SnippetModel type which wraps an sql.DB connection pool
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

	// The SQL statement we want to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	// initialized a new Snippet struct
	var s Snippet

	// Use `row.Scan()` to copy the values from each field in the sql.Row
	// to the corresponding field in the Snippet struct. Arguments to
	// row.Scan are *pointers* to the place the data is copied into.
	// Number of arguments must be exactly the same as the number of
	// columns returned by the statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		// If no rows are returned, then error is returned
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord // custom error.
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

// Return 10 most recent snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {

	stmt := `SELECT id, title, content, created, expires 
	FROM snippets 
	WHERE expires > UTC_TIMESTAMP() 
	ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Closes the resultset once the function returns
	// As long as the resultset is open, it'll keep
	// the underlying database connection open. This
	// can use up all the connections in your pool.
	defer rows.Close()

	// Will hold the Snippet structs returned
	var snippets []Snippet

	// Use `rows.Next()` to iterate throught thr rows in the resultset
	// This prepares the first (and then each subsequent) row to be acted
	// upon by the rows.Scan() method. If iteration over all the rows
	// completes, then the resultset automatically closes itself and frees
	// up the underlying database connection
	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append snippet to slice
		snippets = append(snippets, s)
	}

	// Once the loop has finished, call rows.Err() to
	// retrieve any error encountered during the iteration.
	// ** IMPORTANT ** call this; do not assume a successful
	// iteration over the entire resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
