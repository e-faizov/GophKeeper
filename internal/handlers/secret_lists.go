package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/e-faizov/GophKeeper/internal/models"
)

func (m *SecretsHandlers) GetSecretsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(models.UUIDKey).(string)

	res, err := m.Logic.GetSecretsList(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("SecretsHandlers.GetSecretsList error processing data")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, res)
}
