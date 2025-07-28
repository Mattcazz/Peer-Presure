package post

import (
	"net/http"

	"github.com/Mattcazz/Peer-Presure.git/service/auth"
	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/Mattcazz/Peer-Presure.git/utils"
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
	r.HandleFunc("/post", auth.JWTAuth(h.handleCreatePost)).Methods(http.MethodPost)
	r.HandleFunc("/post", h.handleDeletePost).Methods(http.MethodDelete)
	r.HandleFunc("/:user/posts", auth.JWTAuth(h.handleGetUserPosts)).Methods(http.MethodGet)
	r.HandleFunc("/post/:id", h.handleGetPost).Methods(http.MethodPost)
}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {

	var post types.Post

	err := r.ParseForm()

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	post.Title = r.FormValue("title")
	post.Text = r.FormValue("body")
	if r.FormValue("image") != "" {
		post.ImgURL = r.FormValue("image")
	}
	post.Public = r.FormValue("public") != ""
	post.UserId = r.Context().Value("username").(int)

	err = h.postStore.CreatePost(post)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, post)

}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetPost(w http.ResponseWriter, r *http.Request) {

}
