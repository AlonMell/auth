package httpServer

import (
	"net/http"
	"providerHub/internal/httpServer/api"
	_ "providerHub/internal/httpServer/api"
	"providerHub/pkg/middleware"
)

func Run() error {
	http.Handle("/", middleware.Conveyor(http.HandlerFunc(api.Main), middleware.CORS))
	/*mux := http.NewServeMux()
	mux.HandleFunc("/login/", api.Login)
	mux.HandleFunc("/register/", api.Register)
	mux.HandleFunc("/api/isAdmin/", api.IsAdmin)
	mux.HandleFunc("/", api.Main)*/

	err := http.ListenAndServe(`:8080`, nil)

	return err
}
