package server

import (
	"go-test/internal/server/handlers"
	"go-test/pkg/framework"
	"os"

	"github.com/kpango/glg"
)

type Server struct {
	Router   *framework.Router
	Handlers *handlers.Handlers
}

func (s *Server) Start() {

	app_port := os.Getenv("APP_PORT")
	if app_port == "" {
		app_port = "8080"
	}

	store := framework.NewSessionStore()
	s.Router.Use(framework.SessionMiddleware(store))

	s.Router.Use(framework.CORSMiddleware())
	s.Router.Use(framework.ErrorMiddleware)
	s.Router.Use(LoggerMiddleware)

	s.Router.GET("/people", s.Handlers.GetPeople)
	s.Router.POST("/people", s.Handlers.PostPeople, ValidatePeopleMiddleware)
	s.Router.PUT("/people/:id", s.Handlers.PutPeople, ValidatePeopleMiddleware)
	s.Router.DELETE("/people/:id", s.Handlers.DeletePeople)

	glg.Info("Server listening on port " + app_port)
	s.Router.Run(":" + app_port)
}
