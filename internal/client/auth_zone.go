package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/e-faizov/GophKeeper/internal/cli"
	"github.com/e-faizov/GophKeeper/internal/crypto"
	"github.com/e-faizov/GophKeeper/internal/interfaces"
	"github.com/e-faizov/GophKeeper/internal/models"
	"github.com/e-faizov/GophKeeper/internal/utils"
)

func FirstAction(req interfaces.Requests) {
	for {
		fmt.Println()
		fmt.Println("----------------------")
		fmt.Println("Select action")
		sel := []cli.SelectionItem{
			"New secret",
			"List all secrets",
			"Exit",
		}
		item := cli.Selection(sel)

		switch item {
		case 0:
			newSecret(req)
		case 1:
			list(req)
		case 2:
			os.Exit(0)
		}
	}
}

func list(req interfaces.Requests) {
	secrets, err := req.GetSecretsList()
	if err != nil {
		fmt.Println("Error get secret list", err, "try again")
		return
	}

	if len(secrets) == 0 {
		fmt.Println("Empty")
		return
	}

	names := make([]cli.SelectionItem, 0, len(secrets))

	fmt.Println("Select secret:")
	for _, s := range secrets {
		data, err := crypto.UnCrypt(s.Name)
		if err != nil {
			fmt.Println("Error uncrypt data, try again")
			return
		}

		meta, err := crypto.UnCrypt(s.Meta)
		if err != nil {
			fmt.Println("Error uncrypt data, try again")
			return
		}

		var tpStr string
		switch s.Type {
		case passLoginType:
			tpStr = "(Login-Password)"
		case textType:
			tpStr = "(Text)"
		case bankCardType:
			tpStr = "(Bank card)"
		}

		item := data + " " + tpStr
		if len(meta) != 0 {
			item += "\n" + meta
		}
		names = append(names, cli.SelectionItem(item))
	}
	selectName := cli.Selection(names)

	secret, err := req.GetSecret(secrets[selectName].Uid, secrets[selectName].Version)
	if err != nil {
		fmt.Println("Error get secret", err, "try again")
		return
	}

	err = printSecret(secret)
	if err != nil {
		fmt.Println("Error uncrypt data, try again")
		return
	}

	fmt.Println()
	fmt.Println("----------------------")
	fmt.Println("Select action")
	actionSelection := []cli.SelectionItem{
		"Edit",
		"Remove",
		"Exit",
	}

	selectionResult := cli.Selection(actionSelection)

	switch selectionResult {
	case 0:
		editSecret(req, secret)
	case 1:
		removeSecret(req, secret)
	default:
		return
	}
}

func editSecret(req interfaces.Requests, s models.Secret) {
	var err error
	edit := models.Secret{
		Uid:     s.Uid,
		Version: s.Version,
		Name:    s.Name,
	}
	for {
		data, meta := dataMetaFields(s.Type)

		edit.Data, err = crypto.Crypt(data)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		edit.Meta, err = crypto.Crypt(meta)
		if err != nil {
			fmt.Println("Problem with crypt:", err, "try again")
			continue
		}

		err = req.EditSecret(edit)
		if err != nil {
			fmt.Println(err, "try again")
			continue
		}
		fmt.Println("Done!")
		break
	}
}

func removeSecret(req interfaces.Requests, s models.Secret) {
	err := req.RemoveSecret(s.Uid, s.Version)
	if err != nil {
		fmt.Println("Error when remove secret, try again")
	}
	fmt.Println("Done!")
}

func printSecret(s models.Secret) error {
	unСryptData, err := crypto.UnCrypt(s.Data)
	if err != nil {
		return err
	}

	switch s.Type {
	case passLoginType:
		pl, err := utils.ConvertFromGOB64(unСryptData)
		if err != nil {
			return err
		}
		fmt.Println("Login:", pl.Login)
		fmt.Println("Password:", pl.Password)
	case textType:
		fmt.Println("Text:", unСryptData)
	case bankCardType:
		fmt.Println("Card:", unСryptData)
	default:
		return errors.New("unknown secret type")
	}
	return nil
}
