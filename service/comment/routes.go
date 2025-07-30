package comment

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
	store types.CommentStore
}

func NewHandler(s types.CommentStore) *Handler {

	return &Handler{store: s}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/comment/{id}", h.handleStartComment).Methods(http.MethodPost)
	r.HandleFunc("/post/{id}/comments", h.handleGetPostComments).Methods(http.MethodGet)
	r.HandleFunc("/username/comments", auth.JWTAuth(h.handleGetUserComments)).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/comment", auth.JWTAuth(h.handleCreateComment)).Methods(http.MethodPost)
	r.HandleFunc("/post/{post_id}/comment/{comment_id}", auth.JWTAuth(h.handleDeleteComment)).Methods(http.MethodDelete)
}

func (h *Handler) handleStartComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDStr := vars["id"]

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	web.RenderTemplate(w, "comment-form", map[string]any{"Comment": true, "postID": postID})
}

func (h *Handler) handleGetPostComments(w http.ResponseWriter, r *http.Request) {
	/*	vars := mux.Vars(r)
		postIDStr := vars["id"]

		postID, err := strconv.Atoi(postIDStr)

		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		comments, err := h.store.GetCommentsFromPost(postID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// renderTemplate with comments
	*/
}

func (h *Handler) handleGetUserComments(w http.ResponseWriter, r *http.Request) {
	/*	userID, ok := r.Context().Value(types.CtxKeyUserID).(int)

		if ok {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		comments, err := h.store.GetCommentsFromUser(userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// renderTemplate with comments
	*/
}

func (h *Handler) handleCreateComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	postIDStr := vars["id"]

	postID, err := strconv.Atoi(postIDStr)

	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var comment types.Comment

	comment.PostID = postID

	if r.FormValue("body") == "" {
		http.Error(w, "a comment needs to have text", http.StatusBadRequest)
		return
	}

	comment.Text = r.FormValue("body")

	userID, ok := r.Context().Value(types.CtxKeyUserID).(int)

	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	username, ok := r.Context().Value(types.CtxKeyUsername).(string)

	if !ok {
		http.Error(w, "invalid username", http.StatusBadRequest)
		return
	}
	comment.UserID = userID
	comment.Username = username
	comment.CreatedAt = time.Now()

	_, err = h.store.CreateComment(&comment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// redirect to post comments url
	w.Header().Set("HX-Redirect", fmt.Sprintf("/post/%d", postID))
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleDeleteComment(w http.ResponseWriter, r *http.Request) {

}
