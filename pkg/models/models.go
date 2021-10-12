package models

import ( 
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID int
	Title string
	Cotent string
	Created time.time
	Expires time.time
}