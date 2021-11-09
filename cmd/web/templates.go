package main

import "alexedwards.net/snippetbox/pkg/models"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}