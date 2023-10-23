package framework

import (
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// generateSessionID генерирует новый уникальный ID сессии
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// sessionMiddleware создает новую сессию или восстанавливает существующую
func SessionMiddleware(store *SessionStore) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			// Пытаемся извлечь существующую сессию
			cookie, err := ctx.Request.Cookie("session_id")
			var sessionData map[string]interface{}
			var sessionID string
			if err == nil {
				sessionID = cookie.Value
				sessionData, _ = store.Get(sessionID)
			}

			// Если сессия не существует, создаем новую
			if sessionData == nil {
				sessionData = make(map[string]interface{})
				sessionID = generateSessionID()
				ctx.SetCookie("session_id", sessionID, 3600) // 1 hour
				store.Set(sessionID, sessionData)
			}

			// Добавляем session_id в данные сессии
			sessionData["session_id"] = sessionID

			ctx.SessionData = sessionData

			next(ctx)
		}
	}
}

// GzipMiddleware will apply gzip compression to the response body if the client can accept it
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client can accept the gzip encoding.
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Set the HTTP header to gzip.
		w.Header().Set("Content-Encoding", "gzip")

		// Create a gziped response.
		gz := gzip.NewWriter(w)
		defer gz.Close()

		next.ServeHTTP(GzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (g GzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func StaticCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(allowedIPs ...string) func(next HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if len(allowedIPs) > 0 {
				originIP := strings.Split(ctx.Request.RemoteAddr, ":")[0]
				allowed := false
				for _, ip := range allowedIPs {
					if ip == originIP {
						allowed = true
						break
					}
				}

				if !allowed {
					http.Error(ctx.Writer, "Forbidden", http.StatusForbidden)
					return
				}
			}

			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			if ctx.Request.Method == "OPTIONS" {
				ctx.Writer.WriteHeader(http.StatusOK)
				return
			}

			next(ctx)
		}
	}
}

func ErrorMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("An error occurred: %v", err)
				ctx.JSON(http.StatusInternalServerError, H{"error": fmt.Sprintf("An error occurred: %v", err)})
			}
		}()
		next(ctx)
	}
}
