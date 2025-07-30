package post

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Mattcazz/Peer-Presure.git/service/auth"
	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/Mattcazz/Peer-Presure.git/utils"
	"github.com/Mattcazz/Peer-Presure.git/web"
	"github.com/gorilla/mux"
)

type Handler struct {
	postStore    types.PostStore
	commentStore types.CommentStore
	userStore    types.UserStore
}

func NewHandler(ps types.PostStore, cs types.CommentStore, us types.UserStore) *Handler {
	return &Handler{
		postStore:    ps,
		commentStore: cs,
		userStore:    us,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/post", auth.JWTAuth(h.handleCreatePostGet)).Methods(http.MethodGet)
	r.HandleFunc("/post", auth.JWTAuth(h.handleCreatePostPost)).Methods(http.MethodPost)
	r.HandleFunc("/{username}/posts", auth.JWTAuth(h.handleGetUserPosts)).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}", h.handleGetPost).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/edit", auth.JWTAuth(h.handleEditPostGet)).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}", auth.JWTAuth(h.handleDeletePost)).Methods(http.MethodDelete)
	r.HandleFunc("/post/{id}/edit", auth.JWTAuth(h.handleEditPostPost)).Methods(http.MethodPost)
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
	post.CreatedAt = time.Now()

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

	err = h.commentStore.DeleteCommentsFromPost(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// redirect to getUserPosts
	w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/posts", username))
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	_, err := h.userStore.GetUserByUsername(username)

	if err != nil {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	posts, _ := h.postStore.GetPostsFromUser(username)

	data := types.Data{
		"Posts":    posts,
		"Username": username,
	}

	web.RenderTemplate(w, "user-posts", data)
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

	comments, err := h.commentStore.GetCommentsFromPost(postID)

	if err != nil {
		http.Error(w, "error with the post id", http.StatusBadRequest)
		return
	}

	var userID int
	cookie, err := r.Cookie("auth_token")

	if err != nil {
		userID = 0
	} else {
		tknString := cookie.Value
		tkn, err := auth.ValidateJWTtoken(tknString)

		if err != nil {
			userID = 0
		} else {
			userID, _ = auth.GetUserIdFromJWT(tkn)
		}
	}

	data := types.Data{
		"Post":     post,
		"Comments": comments,
		"UserID":   userID,
	}

	web.RenderTemplate(w, "post-page", data)
}

func (h *Handler) handleEditPostGet(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(types.CtxKeyUserID).(int)

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

	if userID != post.UserId {
		http.Error(w, "You are not allowed to edit this post", http.StatusBadRequest)
		return
	}

	data := types.Data{
		"Post": post,
	}

	web.RenderTemplate(w, "edit-post", data)
}

func (h *Handler) handleEditPostPost(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(types.CtxKeyUserID).(int)

	vars := mux.Vars(r)
	id := vars["id"]

	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postStore.GetPostById(postID)

	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if userID != post.UserId {
		http.Error(w, "You are not allowed to edit this post", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()

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
	post.ID = postID

	err = h.postStore.EditPost(post)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/post/%d", postID))
	w.WriteHeader(http.StatusOK)
}
