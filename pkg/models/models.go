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
	ID int `json:"id"`
	EmpID string `json:"emp_id"`
	EmpName string `json:"emp_name"`
	Role string `json:"role"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}