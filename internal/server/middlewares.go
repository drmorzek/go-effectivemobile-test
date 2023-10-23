package server

import (
	"go-test/pkg/framework"
	"net/http"

	"github.com/kpango/glg"
)

func LoggerMiddleware(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx *framework.Context) {
		defer func() {
			if err := recover(); err != nil {
				glg.Errorf("An error occurred: %v", err)
			}
		}()
		glg.Infof("%v %v", ctx.Request.Method, ctx.Request.URL.Path)
		next(ctx)
	}
}

func ValidatePeopleMiddleware(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx *framework.Context) {
		if err := ctx.ParseJson(); err != nil {
			ctx.JSON(http.StatusBadRequest, framework.H{"error": "Failed to parse JSON"})
			return
		}

		// Perform validation
		if name, ok := ctx.Body["name"].(string); !ok || name == "" {
			glg.Errorf("An error occurred: %v", framework.H{"error": "Missing or invalid 'name' field"})
			ctx.JSON(http.StatusBadRequest, framework.H{"error": "Missing or invalid 'name' field"})
			return
		}
		if surname, ok := ctx.Body["surname"].(string); !ok || surname == "" {
			glg.Errorf("An error occurred: %v", framework.H{"error": "Missing or invalid 'surname' field"})
			ctx.JSON(http.StatusBadRequest, framework.H{"error": "Missing or invalid 'surname' field"})
			return
		}

		if patronymic, ok := ctx.Body["patronymic"].(string); !ok || patronymic == "" {
			ctx.Body["patronymic"] = ""
		}
		next(ctx)
	}
}
