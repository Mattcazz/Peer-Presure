package comment

import (
	"net/http"

	"github.com/Mattcazz/Peer-Presure.git/service/auth"
	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.CommentStore
}

func NewHandler(s types.CommentStore) *Handler {

	return &Handler{store: s}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/post/{id}/comments", h.handleGetPostComments).Methods(http.MethodGet)
	r.HandleFunc("/username/comments", auth.JWTAuth(h.handleGetUserComments)).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/comment", auth.JWTAuth(h.handleCreateCommentPage)).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/comment", auth.JWTAuth(h.handleCreateComment)).Methods(http.MethodPost)
	r.HandleFunc("/post/{id}/comment", auth.JWTAuth(h.handleDeleteComment)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetPostComments(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetUserComments(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleCreateCommentPage(w http.ResponseWriter, r *http.Request) {
	// render create-comment-page
}

func (h *Handler) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	// create comment
}

func (h *Handler) handleDeleteComment(w http.ResponseWriter, r *http.Request) {

}
