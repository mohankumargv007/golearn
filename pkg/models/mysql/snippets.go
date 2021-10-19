package mysql

import (
	"database/sql"
	"alexedwards.net/snippetbox/pkg/models"
	"errors"
)

type SnippetModel struct {
	DB *sql.DB
}

//Inserting Data Into Snippet Model
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	//MySql Statement
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))` 
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err 
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	print(id);
	return 0, nil
}

//get Data Into Snippet Model
func (m *SnippetModel) Get(id int) (*models.Snippet, error) { 
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord 
		} else {
			return nil, err
		}
	}
	//If No Error
	return s, nil
}

//Get Latest Data
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
