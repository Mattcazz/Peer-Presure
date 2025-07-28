package post

import (
	"net/http"
	"strconv"

	"github.com/Mattcazz/Peer-Presure.git/service/auth"
	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/Mattcazz/Peer-Presure.git/utils"
	"github.com/Mattcazz/Peer-Presure.git/web"
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
	r.HandleFunc("/post", auth.JWTAuth(h.handleCreatePostGet)).Methods(http.MethodGet)
	r.HandleFunc("/post", auth.JWTAuth(h.handleCreatePostPost)).Methods(http.MethodPost)
	r.HandleFunc("/post", h.handleDeletePost).Methods(http.MethodDelete)
	r.HandleFunc("/:user/posts", auth.JWTAuth(h.handleGetUserPosts)).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}", h.handleGetPost).Methods(http.MethodGet)
}

func (h *Handler) handleCreatePostPost(w http.ResponseWriter, r *http.Request) {

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
	post.UserId = r.Context().Value(types.CtxKeyUserID).(int)

	p, err := h.postStore.CreatePost(post)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("HX-Redirect", "/post/"+strconv.Itoa(p.ID))
	w.WriteHeader(http.StatusOK)
	// utils.WriteJSON(w, http.StatusOK, &p)
}

func (h *Handler) handleCreatePostGet(w http.ResponseWriter, r *http.Request) {
	web.RenderTemplate(w, "create-post", map[string]any{})
}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postStore.GetPostById(postID)
	if err != nil {
		http.Error(w, "Post doesn't exist", http.StatusBadRequest)
		return
	}

	username := post.UserId

	data := map[string]any{
		"title": post.Title,
		"body": post.Text,
		"username": username,
		"img_url": post.ImgURL,
	}
	web.RenderTemplate(w, "post-page", data)
}
