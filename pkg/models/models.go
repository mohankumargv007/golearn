package models

import ( 
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type Employee struct {
	ID int
	EmpID string
	EmpName string
	Role string
	Created time.Time
	Updated time.Time
}