package main

import "github.com/rhysmah/snippet-box/internal/models"

// Define a templateData type to act as the holding
// structure for any dynamic data that we want to pass
// to our HTML templates. We can only pass one piece
// of dynamic data, so a struct is a way to contain
// one datum composed of many data.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
