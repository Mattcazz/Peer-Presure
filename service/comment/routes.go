package comment

import "github.com/Mattcazz/Peer-Presure.git/types"

type Handler struct {
	store types.CommentStore
}

func NewHandler(s types.CommentStore) *Handler {

	return &Handler{store: s}
}
