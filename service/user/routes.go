package user

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	_, err = auth.JWTAuthWeb(tokenString)

	if err != nil {
		h.handleHomeGuest(w, r)
		return
	}

	h.handleHomeUser(w, r)
}

func (h *Handler) handleHomeGuest(w http.ResponseWriter, r *http.Request) {
	web.RenderTemplate(w, "home", map[string]any{})
}

func (h *Handler) handleHomeUser(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postStore.GetLastTenPosts()

	if err != nil {
		http.Error(w, "WTF!!", http.StatusBadRequest)
	}

	d := types.Data{
		"Posts": posts,
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

	secret := os.Getenv("JWT_SECRET")

	token, err := auth.CreateJWT([]byte(secret), u.ID, u.UserName)

	if err != nil {
		web.RenderTemplate(w, "login-form", map[string]any{"Error": "Unauthorized"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	})

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	web.RenderTemplate(w, "login", map[string]any{})
}

func (h *Handler) handleRegisterPost(w http.ResponseWriter, r *http.Request) {

	var payload types.RegisterUserPayload
	err := r.ParseForm()
	// err := utils.ParseJSON(r, &payload)

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

	err = h.userStore.CreateUser(types.User{
		UserName:  payload.UserName,
		Email:     payload.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	})

	if err != nil {
		web.RenderTemplate(w, "register-form", map[string]any{"Error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
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

	utils.WriteJSON(w, http.StatusOK, nil)
}
