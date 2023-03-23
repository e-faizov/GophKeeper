package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

type ClientConfig struct {
	Address string `env:"ADDRESS"`
}

var (
	cliCfg       ClientConfig
	cliCfgInited bool
)

func GetClientConfig() ClientConfig {
	if !cliCfgInited {
		flag.StringVar(&(cliCfg.Address), "a", "http://localhost:8080", "ADDRESS")

		flag.Parse()
		if err := env.Parse(&cliCfg); err != nil {
			panic(err)
		}

		cliCfgInited = true
	}
	return cliCfg
}
