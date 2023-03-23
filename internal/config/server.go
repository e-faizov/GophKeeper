package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

type ServerConfig struct {
	Address     string `env:"ADDRESS"`
	DatabaseDsn string `env:"DATABASE_DSN"`
}

var (
	serverCfg    ServerConfig
	srvCfgInited bool
)

func GetServerConfig() ServerConfig {
	if !srvCfgInited {
		flag.StringVar(&serverCfg.Address, "a", "localhost:8080", "ADDRESS")
		flag.StringVar(&serverCfg.DatabaseDsn, "d", "", "KEY")

		flag.Parse()
		if err := env.Parse(&serverCfg); err != nil {
			panic(err)
		}
		srvCfgInited = true
	}
	return serverCfg
}
