package Middleware

import (
	"net/http"
	"github.com/go-chi/cors"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "APIKey"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(next)
}
