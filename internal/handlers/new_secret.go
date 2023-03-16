package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/e-faizov/GophKeeper/internal/interfaces"
	"github.com/e-faizov/GophKeeper/internal/models"
)

type SecretsHandlers struct {
	Logic interfaces.CryptoLogic
}

func (m *SecretsHandlers) NewSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(models.UUIDKey).(string)

	sec, err := unmarshalSecret(r)
	if err != nil {
		log.Error().Err(err).Msg("SecretsHandlers.NewSecret error parse body")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = m.Logic.NewSecret(ctx, userID, sec)
	if err != nil {
		log.Error().Err(err).Msg("SecretsHandlers.NewSecret error processing data")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func unmarshalSecret(r *http.Request) (models.Secret, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.Secret{}, err
	}

	var data models.Secret
	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.Secret{}, err
	}
	return data, nil
}
