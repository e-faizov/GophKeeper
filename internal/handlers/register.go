package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/rs/zerolog/log"

	"github.com/e-faizov/GophKeeper/internal/interfaces"
	"github.com/e-faizov/GophKeeper/internal/models"
)

type PersonsHandlers struct {
	Store     interfaces.UserStorage
	TokenAuth *jwtauth.JWTAuth
}

func (p *PersonsHandlers) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := unmarshalUser(r)
	if err != nil {
		log.Error().Err(err).Msg("User.Register error unmarshal data")
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	ok, uid, err := p.Store.Register(ctx, user.Login, user.Password)
	if err != nil {
		log.Error().Err(err).Msg("User.Register sql error")
		http.Error(w, "wrong body", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "already created", http.StatusConflict)
		return
	}

	token, err := p.token(uid)
	if err != nil {
		log.Error().Err(err).Msg("User.Register error create token")
		p.Logout(w, r)
		http.Error(w, "wrong body", http.StatusInternalServerError)
		return
	}

	setCookie(w, token)
	w.WriteHeader(http.StatusCreated)
}

func (p *PersonsHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    "",
	})
}

func (p *PersonsHandlers) token(user string) (string, error) {
	_, tokenString, err := p.TokenAuth.Encode(map[string]interface{}{models.UserUUID: user})
	return tokenString, err
}

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    token,
	})
}

func unmarshalUser(r *http.Request) (models.User, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.User{}, err
	}

	var data models.User
	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.User{}, err
	}
	return data, nil
}
