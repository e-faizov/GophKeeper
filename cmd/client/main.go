package main

import (
	"github.com/e-faizov/GophKeeper/internal/client"
	"github.com/e-faizov/GophKeeper/internal/config"
	"github.com/e-faizov/GophKeeper/internal/interfaces"
	"github.com/e-faizov/GophKeeper/internal/network"
)

func main() {
	cfg := config.GetClientConfig()
	var req interfaces.Requests
	req = &network.HttpsRequests{
		Url: cfg.Address,
	}

	//unauth zone
	client.RegOrAuth(req)

	//req.Auth("asdfg", "asdfg")
	//crypto.SetPasswors("asdfg")

	//auth zone
	client.FirstAction(req)
}
