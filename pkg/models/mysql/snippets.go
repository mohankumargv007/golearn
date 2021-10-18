package mysql

import (
	"database/sql"
	//Importing Models
	"alexedwards.net/snippetbox/pkg/models"
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
	return 0, nil
}

//get Data Into Snippet Model
func (m *SnippetModel) Get(id int) (*models.Snippet, error) { 
	return nil, nil
}

//Get Latest Data
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
