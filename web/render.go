package web

import (
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func LoadTemplates() {
	var err error
	templates, err = template.ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func RenderTemplate(w http.ResponseWriter, templateName string, data map[string]any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Printf("Template execution error (%s): %v", templateName, err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

/*
func RenderTemplate2(w http.ResponseWriter, templateName string, title string, data map[string]any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	contentPath := filepath.Join("web/templates", templateName+".html")
	bodyBytes, err := os.ReadFile(contentPath)

	if err != nil {
		log.Printf("Failed to read template: %v", err)
		http.Error(w, "Page not found", http.StatusInternalServerError)
		return
	}

	//	contentPath = filepath.Join("web/templates", templateName+".css")
	//	styleBytes, err := os.ReadFile(contentPath)

	//	if err != nil {
	//		log.Printf("Failed to read template: %v", err)
	//		http.Error(w, "Page not found", http.StatusInternalServerError)
	//		return
	//	}

	data["Title"] = title
	data["Body"] = template.HTML(bodyBytes)

	err = baseTemplate.Execute(w, data)
	if err != nil {
		log.Printf("renderTemplate error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
*/
