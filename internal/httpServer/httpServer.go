package httpServer

import (
	"net/http"
	"providerHub/internal/httpServer/endpoint"
	_ "providerHub/internal/httpServer/endpoint"
	"providerHub/pkg/middleware"
)

func Run() error {
	http.Handle("/", middleware.Conveyor(http.HandlerFunc(endpoint.Main), middleware.CORS))
	/*mux := http.NewServeMux()
	mux.HandleFunc("/login/", api.Login)
	mux.HandleFunc("/register/", api.Register)
	mux.HandleFunc("/api/isAdmin/", api.IsAdmin)
	mux.HandleFunc("/", api.Main)*/

	err := http.ListenAndServe(`:8080`, nil)

	return err
}
