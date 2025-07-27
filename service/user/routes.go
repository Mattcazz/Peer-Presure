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
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(s types.UserStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleLogin called")
	var payload types.LoginUserPayload
	err := r.ParseForm()
	// err := utils.ParseJSON(r, &payload)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	payload.Email = r.FormValue("email")
	payload.Password = r.FormValue("password")

	log.Printf("Got POST request with email '%v' and password '%v'\n", payload.Email, payload.Password)

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", errors))
	}

	u, err := h.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ValidatePassword(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := os.Getenv("JWT_SECRET")

	token, err := auth.CreateJWT([]byte(secret), u.ID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {

	var payload types.RegisterUserPayload
	err := r.ParseForm()
	// err := utils.ParseJSON(r, &payload)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	payload.Email = r.FormValue("email")
	payload.UserName = r.FormValue("username")
	payload.Password = r.FormValue("password")

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", errors))
		return
	}

	_, err = h.store.GetUserByEmail(payload.Email)

	if err == nil { // err == nil that means that it found a user with the given email
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		UserName:  payload.UserName,
		Email:     payload.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
