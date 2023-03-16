package main

import (
	"github.com/e-faizov/GophKeeper/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	var srv server.CryptoServer

	err := srv.StartServer()
	if err != nil {
		log.Error().Err(err).Msg("can't start server")
	}
}
