package main

import (
	"github.com/rs/zerolog/log"

	"github.com/e-faizov/GophKeeper/internal/config"
	"github.com/e-faizov/GophKeeper/internal/server"
)

func main() {
	cfg := config.GetServerConfig()

	var srv server.CryptoServer

	err := srv.StartServer(cfg)
	if err != nil {
		log.Error().Err(err).Msg("can't start server")
	}
}
