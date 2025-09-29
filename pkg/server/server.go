package server

import (
	"fmt"
	"net/http"
	"os"
)

func Run() error {
	webDir := "./web"
	port := "7540"

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}
	addr := ":" + port

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	return nil
}
