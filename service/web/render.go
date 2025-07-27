package web

import (
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func loadTemplates(pattern string) {
	var err error
	templates, err = template.ParseGlob(pattern)
	if (err != nil) {
		log.Fatalf("failed to parse templates: %v", err)
	}
}

func renderTemplate(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, name, data)
	if (err != nil) {
		log.Printf("renderTemplate error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
