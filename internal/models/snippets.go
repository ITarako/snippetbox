package models

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type SnippetModel struct {
	DB *sqlx.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	var id int64

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, NOW(), NOW() + ($3 || ' days')::interval)
	RETURNING id`

	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() AND id = $1`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
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
