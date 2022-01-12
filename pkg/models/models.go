package models

import (
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type Employee struct {
	ID      int       `json:"id"`
	EmpID   string    `json:"emp_id"`
	EmpName string    `json:"emp_name"`
	Role    string    `json:"role"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

var (
	ErrNoRecord           = errors.New("models:  no founds")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
