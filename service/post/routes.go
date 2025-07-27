package post

import (
	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/gorilla/mux"
)

type Handler struct {
	postStore    types.PostStore
	commentStore types.CommentStore
}

func NewHandler(ps types.PostStore, cs types.CommentStore) *Handler {
	return &Handler{
		postStore:    ps,
		commentStore: cs,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {

}
