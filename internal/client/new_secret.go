package client

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/e-faizov/GophKeeper/internal/cli"
	"github.com/e-faizov/GophKeeper/internal/crypto"
	"github.com/e-faizov/GophKeeper/internal/interfaces"
	"github.com/e-faizov/GophKeeper/internal/models"
	"github.com/e-faizov/GophKeeper/internal/utils"
)

const (
	secretNameLabel = "Enter secret name"
	metaLabel       = "Enter meta"
)

const (
	passLoginType = 0
	textType      = 1
	bankCardType  = 2
)

func newSecret(req interfaces.Requests) {
	var err error
	for {
		fmt.Println("Select secret type")
		typeSel := []cli.SelectionItem{
			"Login/Password",
			"Text",
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

		data, secretMeta := dataMetaFields(secretType)

		var sec models.Secret

		sec.Type = secretType

		sec.Name, err = crypto.Crypt(secretName)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		sec.Data, err = crypto.Crypt(data)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		sec.Meta, err = crypto.Crypt(secretMeta)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		err = req.NewSecret(sec)
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
	return utils.ConvertToGOB64(lp)
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
