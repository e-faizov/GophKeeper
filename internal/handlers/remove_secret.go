package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/e-faizov/GophKeeper/internal/models"
)

func (m *SecretsHandlers) RemoveSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(models.UUIDKey).(string)

	sec, err := unmarshalSecret(r)
	if err != nil {
		log.Error().Err(err).Msg("SecretsHandlers.EditSecret error parse body")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = m.Logic.RemoveSecret(ctx, userID, sec)
	if err != nil {
		log.Error().Err(err).Msg("SecretsHandlers.EditSecret error processing data")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
