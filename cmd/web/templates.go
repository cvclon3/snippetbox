package main

import (
	"html/template"
	"path/filepath"
	"time"
	"io/fs"

	"snippetbox.cvclon3.net/internal/models"
	"snippetbox.cvclon3.net/ui"
)


type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken string
}


func humanDate(t time.Time) string {
	// Why we use this date:
	// https://go.dev/src/time/format.go
	return t.Format("02 Jan 2006 at 15:04")
}


var functions = template.FuncMap{
	"humanDate": humanDate,
}


func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*tmpl",
			page,
		}

		tmpl, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}