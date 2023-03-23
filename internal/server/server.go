package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"

	"github.com/e-faizov/GophKeeper/internal/config"
	"github.com/e-faizov/GophKeeper/internal/handlers"
	"github.com/e-faizov/GophKeeper/internal/logic"
	"github.com/e-faizov/GophKeeper/internal/middlewares"
	"github.com/e-faizov/GophKeeper/internal/stores"
)

var tokenAuth *jwtauth.JWTAuth

type CryptoServer struct {
	srv http.Server
}

// StartServer - функция запуска сервера
func (s *CryptoServer) StartServer(cfg config.ServerConfig) error {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	pg, err := stores.NewPgStore(cfg.DatabaseDsn)
	if err != nil {
		return err
	}

	cl := logic.CryptoLogicImpl{
		Store: pg,
	}

	persons := handlers.PersonsHandlers{
		Store:     pg,
		TokenAuth: tokenAuth,
	}
	secrets := handlers.SecretsHandlers{
		Logic: &cl,
	}

	r := chi.NewRouter()
	//r.Use(middleware.Compress(5))

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
				r.Get("/getAll", secrets.GetSecretsList)
			})
	})

	return http.ListenAndServe(cfg.Address, r)
}
