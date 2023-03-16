package client

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/e-faizov/GophKeeper/internal/cli"
	"github.com/e-faizov/GophKeeper/internal/crypto"
	"github.com/e-faizov/GophKeeper/internal/models"
	"github.com/e-faizov/GophKeeper/internal/network"
	"github.com/e-faizov/GophKeeper/internal/utils"
)

const (
	secretNameLabel = "Enter secret name"
	metaLabel       = "Enter meta"
)

const (
	passLoginType = 0
	textType      = 1
	binaryType    = 2
	bankCardType  = 3
)

func newSecret() {
	var err error
	for {
		fmt.Println("Select secret type")
		typeSel := []cli.SelectionItem{
			"Login/Password",
			"Text",
			"Binary",
			"Bank card",
		}
		secretType := cli.Selection(typeSel)

		var secretName string
		nameField := []cli.EnterDataItem{
			{
				Name:   secretNameLabel,
				Data:   &secretName,
				Verify: cli.NotEmpty,
			},
		}
		cli.EnterData(nameField)

		var data string
		switch secretType {
		case passLoginType:
			data, err = loginPassPage()
			if err != nil {
				fmt.Println(err, "try again")
				continue
			}
		case textType:
			fields := []cli.EnterDataItem{
				{
					Name:   "Enter text",
					Data:   &data,
					Verify: cli.NotEmpty,
				},
			}
			cli.EnterData(fields)
		case binaryType:
			data, err = binaryPage()
			if err != nil {
				fmt.Println(err, "try again")
				continue
			}
		case bankCardType:
			fields := []cli.EnterDataItem{
				{
					Name:   "Enter bank card",
					Data:   &data,
					Verify: goluhn.Validate,
				},
			}
			cli.EnterData(fields)
		}

		var secretMeta string
		metaField := []cli.EnterDataItem{
			{
				Name: metaLabel,
				Data: &secretMeta,
			},
		}
		cli.EnterData(metaField)

		var sec models.Secret

		sec.Type = secretType

		sec.Data1, err = crypto.Crypt(secretName)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		sec.Data2, err = crypto.Crypt(data)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		sec.Data3, err = crypto.Crypt(secretMeta)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		err = network.NewSecret(sec)
		if err != nil {
			fmt.Println("Problem with sending secret to server:", err, "try again")
			continue
		}

		fmt.Println("Done!")
		break
	}
}

func loginPassPage() (string, error) {
	var lp models.User
	fields := []cli.EnterDataItem{
		{
			Name:   "Login",
			Data:   &lp.Login,
			Verify: cli.NotEmpty,
		},
		{
			Name:   "Password",
			Data:   &lp.Password,
			Verify: cli.NotEmpty,
		},
	}
	cli.EnterData(fields)
	return utils.ToGOB64(lp)
}

func binaryPage() (string, error) {
	var fileName string
	fields := []cli.EnterDataItem{
		{
			Name:   "Enter file path with binary data",
			Data:   &fileName,
			Verify: cli.NotEmpty,
		},
	}
	cli.EnterData(fields)
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
