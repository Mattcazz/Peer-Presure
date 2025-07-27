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
	loadTemplates()
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", h.handleLoginPage).Methods("GET")
	r.HandleFunc("/register", h.handleRegisterPage).Methods("GET")
}

func (h *Handler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", "Login", map[string]any{})
}

func (h *Handler) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", "Register", map[string]any{})
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
