package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (p *PersonsHandlers) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := unmarshalUser(r)
	if err != nil {
		log.Error().Err(err).Msg("User.Login error unmarshal data")
		p.Logout(w, r)
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	ok, uid, err := p.Store.Login(ctx, user.Login, user.Password)
	if err != nil {
		log.Error().Err(err).Msg("User.Login error verify user")
		p.Logout(w, r)
		http.Error(w, "wrong body", http.StatusInternalServerError)
		return
	}

	if !ok {
		p.Logout(w, r)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := p.token(uid)
	if err != nil {
		log.Error().Err(err).Msg("User.Login error create token")
		p.Logout(w, r)
		http.Error(w, "wrong body", http.StatusInternalServerError)
		return
	}

	setCookie(w, token)
}
