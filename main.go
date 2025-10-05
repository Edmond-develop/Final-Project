package main

import (
	"fmt"
	"go1f/pkg/api"
	"go1f/pkg/db"
	"go1f/pkg/server"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	err := db.Init(dbFile)
	if err != nil {
		fmt.Println("init scheduler error:", err)
		return
	}
	defer db.DB.Close()
	api.Init()
	err = server.Run()
	if err != nil {
		fmt.Println("server run error:", err)
		return
	}
}
