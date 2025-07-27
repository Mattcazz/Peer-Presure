package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var baseTemplate *template.Template

func loadTemplates() {
	var err error
	baseTemplate, err = template.ParseFiles("web/templates/base.html")
	if (err != nil) {
		log.Fatalf("failed to parse base template: %v", err)
	}
}

func renderTemplate(w http.ResponseWriter, templateName string, title string, data map[string]any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	contentPath := filepath.Join("web/templates", templateName + ".html")
	bodyBytes, err := os.ReadFile(contentPath)

	if (err != nil) {
		log.Printf("Failed to read template: %v", err)
		http.Error(w, "Page not found", http.StatusInternalServerError)
		return
	}

	contentPath = filepath.Join("web/templates", templateName + ".css")
	styleBytes, err := os.ReadFile(contentPath)

	if (err != nil) {
		log.Printf("Failed to read template: %v", err)
		http.Error(w, "Page not found", http.StatusInternalServerError)
		return
	}

	data["Title"] = title
	data["Style"] = template.CSS(styleBytes)
	data["Body"] = template.HTML(bodyBytes)

	err = baseTemplate.Execute(w, data)
	if (err != nil) {
		log.Printf("renderTemplate error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
