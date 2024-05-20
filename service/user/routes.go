package user

import (
	"fmt"
	"net/http"

	"github.com/Maliud/Backend-API-in-Golang/config"
	"github.com/Maliud/Backend-API-in-Golang/service/auth"
	"github.com/Maliud/Backend-API-in-Golang/types"
	"github.com/Maliud/Backend-API-in-Golang/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POS")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// JSON'nı al
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// payload doğrulama

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Geçersiz Payload &v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Geçersiz email veya parola"))
		return
	}
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Geçersiz Parola"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	} 
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token:": token})

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// JSON'nı al
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// payload doğrulama

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Geçersiz Payload &v", errors))
		return
	}

	// kullanıcının var olup olmadığını kontrol et

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("e-posta %s kullanıcısı zaten var", payload.Email))
		return
	}

	hashedPassword, err := auth.HashedPassword(payload.Password)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// eğer değilse yeni kullanıcıyı yaratırız

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
