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
