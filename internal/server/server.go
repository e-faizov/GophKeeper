package server

import (
	"net/http"

	"github.com/e-faizov/GophKeeper/internal/handlers"
	"github.com/e-faizov/GophKeeper/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

type CryptoServer struct {
	srv http.Server
}

// StartServer - функция запуска сервера
func (s *CryptoServer) StartServer() error {

	persons := handlers.PersonsHandlers{}
	secrets := handlers.SecretsHandlers{}

	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", persons.Register)
			r.Post("/login", persons.Auth)
		})

		r.With(jwtauth.Verifier(tokenAuth), middlewares.Auth).
			Route("/secrets", func(r chi.Router) {
				r.Post("/new", secrets.NewSecret)
				r.Post("/edit", secrets.EditSecret)
				r.Post("/remove", secrets.RemoveSecret)
				r.Post("/get", secrets.GetSecret)
				r.Post("/getAll", secrets.GetSecretsList)
			})
	})

	return http.ListenAndServe(":8080", r)
}
