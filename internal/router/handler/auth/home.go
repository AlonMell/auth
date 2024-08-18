package auth

import (
	"encoding/json"
	"net/http"
)

type SimpleResponse struct {
	Message string `json:"message"`
	Owner   string `json:"owner"`
}

func Main(w http.ResponseWriter, r *http.Request) {
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
