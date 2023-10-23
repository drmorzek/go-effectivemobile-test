package main

import (
	"go-test/internal/db"
	"go-test/internal/server"
	"go-test/internal/server/handlers"
	"go-test/pkg/framework"

	"github.com/joho/godotenv"
	"github.com/kpango/glg"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		glg.Errorf("Error loading .env file")
	}

	DB := db.ConnectDB()
	defer DB.Close()

	router := framework.NewRouter()

	handlers := &handlers.Handlers{
		DB: DB,
	}

	server := &server.Server{
		Router:   router,
		Handlers: handlers,
	}

	server.Start()
}
