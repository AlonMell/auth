package api

import (
	"encoding/json"
	"net/http"
	"providerHub/internal/domain/models"
)

type SimpleResponse struct {
	Message string `json:"message"`
	Owner   string `json:"owner"`
}

func Main(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	resp := SimpleResponse{"Hello", "ProviderHub"}
	res, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}

// in: login, password
// out: userId

// Сделать DTO для request, response
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	//Hash пароля
	//Положить User в бд
}

// in: login, password,
// out: token

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	//Проверить есть ли такой пользовать, забрать hash пароля
	//Сравнить пришедший хэш пароля с хранящимся в бд
	//Сгенерировать токен и вырнуть его в куки
}

// in: userId
// out: isAdmin

func IsAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	res, err := json.Marshal(struct {
		IsAdmin bool `json:"isAdmin"`
	}{IsAdmin: false})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(res)
}
