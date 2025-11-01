package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"scvis.vaccsedu.org/internal/models"
	"scvis.vaccsedu.org/ui"
)

type templateData struct {
	CurrentYear      int
	Project          models.Project
	Projects         []models.Project
	Issue            *models.Issue
	Issues           []*models.Issue
	Form             any
	Flash            string
	IsAuthenticated  bool
	IsAdministrator  bool
	CSRFToken        string
	CFile            models.CFile
	InputFile        models.InputFile
	GitHubConfigured bool
	GitHubRepo       string
	CSPNonce         string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"split":     strings.Split,
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
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
