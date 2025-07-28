package post

import (
	"fmt"
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
	r.HandleFunc("/{user}/posts", auth.JWTAuth(h.handleGetUserPosts)).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}", h.handleGetPost).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}", auth.JWTAuth(h.handleDeletePost)).Methods(http.MethodDelete)
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
	post.Username = r.Context().Value(types.CtxKeyUsername).(string)

	p, err := h.postStore.CreatePost(post)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("HX-Redirect", "/post/"+strconv.Itoa(p.ID))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleCreatePostGet(w http.ResponseWriter, r *http.Request) {
	web.RenderTemplate(w, "create-post", map[string]any{})
}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error: invalid post ID"))
		return
	}

	post, err := h.postStore.GetPostById(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error: post does not exist"))
		return
	}

	userID := r.Context().Value(types.CtxKeyUserID).(int)

	username, ok := r.Context().Value(types.CtxKeyUsername).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: username missing"))
		return
	}

	if userID != post.UserId {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unauthorized"))
		return
	}

	err = h.postStore.DeletePost(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// redirect to getUserPosts
	w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/posts", username))
	w.WriteHeader(http.StatusNoContent) // Or you can return a success snippet instead
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

	if err != nil {
		http.Error(w, "user does not exists (check userID)", http.StatusBadRequest)
		return
	}

	data := map[string]any{
		"title":    post.Title,
		"body":     post.Text,
		"username": post.Username,
		"img_url":  post.ImgURL,
		"postID":   post.ID,
	}
	web.RenderTemplate(w, "post-page", data)
}
