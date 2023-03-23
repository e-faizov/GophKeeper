package client

import (
	"fmt"
	"os"

	"github.com/e-faizov/GophKeeper/internal/cli"
	"github.com/e-faizov/GophKeeper/internal/crypto"
	"github.com/e-faizov/GophKeeper/internal/interfaces"
)

func getJwt(req func(string, string) error) error {
	var login, pass string
	items := []cli.EnterDataItem{
		{
			Name:   "Login",
			Data:   &login,
			Verify: cli.MoreThan(5),
		},
		{
			Name:   "Password",
			Data:   &pass,
			Verify: cli.MoreThan(5),
		},
	}

	cli.EnterData(items)
	crypto.SetPasswors(pass)

	authPass, err := crypto.Hmac(login)
	if err != nil {
		return err
	}

	err = req(login, authPass)

	if err != nil {
		return err
	}
	return nil
}

func RegOrAuth(req interfaces.Requests) {
LOOP:
	for {
		ra := []cli.SelectionItem{
			"Registration",
			"Auth",
			"Exit",
		}

		itemSel := cli.Selection(ra)

		switch itemSel {
		case 0:
			err := getJwt(req.Registration)
			if err != nil {
				fmt.Println(err, "try again")
				continue
			}
			break LOOP
		case 1:
			err := getJwt(req.Auth)
			if err != nil {
				fmt.Println(err, "try again")
				continue
			}
			break LOOP
		case 2:
			os.Exit(0)
			return
		}
	}
}
