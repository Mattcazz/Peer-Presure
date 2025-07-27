package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	templates *template.Template
}

func NewHandler() *Handler {
	loadTemplates("web/templates/*.html")
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", h.handleLoginPage).Methods("GET")
}

func (h *Handler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", map[string]any {
		"title": "Login Page",
	})
}

/*
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
		"title": "Login Page",
	})
	if err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
*/
