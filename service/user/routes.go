package user

import (
	"fmt"
	"log"
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
	userStore types.UserStore
	postStore types.PostStore
}

func NewHandler(us types.UserStore, ps types.PostStore) *Handler {

	return &Handler{
		userStore: us,
		postStore: ps}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", h.handleHome).Methods(http.MethodGet)
	router.HandleFunc("/home/{username}", auth.JWTAuth(h.handleHomeUser)).Methods(http.MethodGet)
	router.HandleFunc("/{username}/friends", auth.JWTAuth(h.handleUserFriends)).Methods(http.MethodGet)
	router.HandleFunc("/friends/delete", auth.JWTAuth(h.handleDeleteFriend)).Methods(http.MethodPost)
	router.HandleFunc("/login", h.handleLoginPost).Methods(http.MethodPost)
	router.HandleFunc("/login", h.handleLoginGet).Methods(http.MethodGet)
	router.HandleFunc("/register", h.handleRegisterPost).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegisterGet).Methods(http.MethodGet)
	router.HandleFunc("/logout", h.handleLogout).Methods(http.MethodPost)
}

func (h *Handler) handleHome(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		h.handleHomeGuest(w, r)
		return
	}

	tokenString := cookie.Value

	token, err := auth.JWTAuthWeb(tokenString)

	if err != nil {
		h.handleHomeGuest(w, r)
		return
	}

	username, err := auth.GetUsernameFromJWT(token)

	if err != nil {
		h.handleHomeGuest(w, r)
		return
	}

	_, err = h.userStore.GetUserByUsername(username)

	if err != nil {
		h.handleHomeGuest(w, r)
		return
	}

	url := fmt.Sprintf("/home/%s", username)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) handleHomeGuest(w http.ResponseWriter, r *http.Request) {
	web.RenderTemplate(w, "home", map[string]any{})
}

func (h *Handler) handleHomeUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]

	pageNumber := 1
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			pageNumber = p
		}
	}

	u, err := h.userStore.GetUserByUsername(username)

	if err != nil {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := utils.GetUserIdFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userId != u.ID {
		http.Error(w, "Unauthorized: you are not this user", http.StatusBadRequest)
		return
	}

	posts, pagination, err := paginateFeed(h, userId, pageNumber, username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d := types.Data{
		"Posts":      posts,
		"Username":   username,
		"Pagination": pagination,
	}

	web.RenderTemplate(w, "feed", d)
}

func (h *Handler) handleLoginPost(w http.ResponseWriter, r *http.Request) {

	var payload types.LoginUserPayload
	err := r.ParseForm()

	if err != nil {
		web.RenderTemplate(w, "login-form", map[string]any{"Error": err.Error()})
		return
	}

	payload.Email = r.FormValue("email")
	payload.Password = r.FormValue("password")

	if err := utils.Validate.Struct(payload); err != nil {
		web.RenderTemplate(w, "login-form", map[string]any{"Error": "Invalid data. Please fill all fields"})
		return
	}

	u, err := h.userStore.GetUserByEmail(payload.Email)

	if err != nil {
		web.RenderTemplate(w, "login-form", map[string]any{"Error": "Not found, invalid email or password"})
		return
	}

	if !auth.ValidatePassword(u.Password, []byte(payload.Password)) {
		web.RenderTemplate(w, "login-form", map[string]any{"Error": "Not found, invalid email or password"})
		return
	}

	err = loginUser(w, u.ID, u.UserName)

	if err != nil {
		web.RenderTemplate(w, "login-form", map[string]any{"Error": err.Error()})
		return
	}

	h.handleHomeUser(w, r)
}

func (h *Handler) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	web.RenderTemplate(w, "login", map[string]any{})
}

func (h *Handler) handleRegisterPost(w http.ResponseWriter, r *http.Request) {

	var payload types.RegisterUserPayload
	err := r.ParseForm()

	if err != nil {
		web.RenderTemplate(w, "register-form", map[string]any{"Error": err.Error()})
	}

	payload.Email = r.FormValue("email")
	payload.UserName = r.FormValue("username")
	payload.Password = r.FormValue("password")

	log.Println(payload)

	if err := utils.Validate.Struct(payload); err != nil {
		web.RenderTemplate(w, "register-form", map[string]any{"Error": "Invalid data. Please fill all fields"})

		return
	}

	_, err = h.userStore.GetUserByEmail(payload.Email)

	if err == nil { // err == nil that means that it found a user with the given email
		web.RenderTemplate(w, "register-form", map[string]any{"Error": fmt.Sprintf("User with email %s already exists", payload.Email)})
		return
	}

	_, err = h.userStore.GetUserByUsername(payload.UserName)

	if err == nil { // err == nil that means that it found a user with the given email
		web.RenderTemplate(w, "register-form", map[string]any{"Error": fmt.Sprintf("Username %s already exists", payload.UserName)})
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		web.RenderTemplate(w, "register-form", map[string]any{"Error": err.Error()})
		return
	}

	user := &types.User{
		UserName:      payload.UserName,
		Email:         payload.Email,
		Password:      hashedPassword,
		ProfilePicUrl: types.AvatarURL,
		CreatedAt:     time.Now(),
	}

	err = h.userStore.CreateUser(user)

	if err != nil {
		web.RenderTemplate(w, "register-form", map[string]any{"Error": err.Error()})
		return
	}

	loginUser(w, user.ID, user.UserName)

	w.Header().Set("HX-Redirect", fmt.Sprintf("/home/%s", user.UserName))
	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) handleRegisterGet(w http.ResponseWriter, r *http.Request) {

	web.RenderTemplate(w, "register", map[string]any{})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {

	// we override the cookie with an expiring one so that it is no longer working

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HttpOnly: true,
		Secure:   false,
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleUserFriends(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]

	user, err := h.userStore.GetUserByUsername(username)

	if err != nil {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	friends, err := h.userStore.GetUserFriends(user.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := types.Data{
		"Users":    friends,
		"Username": username,
	}

	web.RenderTemplate(w, "user-friends", data)
}

func (h *Handler) handleDeleteFriend(w http.ResponseWriter, r *http.Request) {
	currentUserId, err := utils.GetUserIdFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	otherUserIdStr := r.FormValue("friend_id")

	otherUserId, err := strconv.Atoi(otherUserIdStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userStore.DeleteFriend(currentUserId, otherUserId)

	if err != nil {
		http.Error(w, "This 2 users are not friends", http.StatusBadRequest)
		return
	}

	username, err := utils.GetUsernameFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/friends", username))
	w.WriteHeader(http.StatusAccepted)
}
