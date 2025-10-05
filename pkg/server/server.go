package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

var Router *chi.Mux

func Run() error {
	InitRouter()
	webDir := "./web"
	port := "7540"

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}
	addr := ":" + port
	fileWeb := http.FileServer(http.Dir(webDir))
	Router.Handle("/*", fileWeb)
	err := http.ListenAndServe(addr, Router)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	return nil
}
func InitRouter() {
	if Router == nil {
		Router = chi.NewRouter()
	}
}
