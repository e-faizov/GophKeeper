package client

import (
	"fmt"
	"os"

	"github.com/e-faizov/GophKeeper/internal/cli"
	"github.com/e-faizov/GophKeeper/internal/crypto"
	"github.com/e-faizov/GophKeeper/internal/models"
	"github.com/e-faizov/GophKeeper/internal/network"
)

func FirstAction() {
	for {
		fmt.Println("Select action")
		sel := []cli.SelectionItem{
			"New secret",
			"List all secrets",
			"Exit",
		}
		item := cli.Selection(sel)

		switch item {
		case 0:
			newSecret()
		case 1:
			list()
		case 2:
			os.Exit(0)
		}
	}
}

func list() {
	secrets, err := network.GetSecretsList()
	if err != nil {
		fmt.Println("Error get secret list", err, "try again")
		return
	}

	names := make([]cli.SelectionItem, 0, len(secrets))

	for _, s := range secrets {
		data, err := crypto.UnCrypt(s.Data1)
		if err != nil {
			fmt.Println("Error uncrypt data, try again")
			return
		}
		names = append(names, cli.SelectionItem(data))
	}
	selectName := cli.Selection(names)

	secret, err := network.GetSecret(secrets[selectName].Uid, secrets[selectName].Version)
	if err != nil {
		fmt.Println("Error get secret", err, "try again")
		return
	}

	err = printSecret(secret)
	if err != nil {
		fmt.Println("Error uncrypt data, try again")
		return
	}

	actionSelection := []cli.SelectionItem{
		"Edit",
		"Remove",
		"Exit",
	}

	selectionResult := cli.Selection(actionSelection)

	switch selectionResult {
	case 0:
		editSecret(secret)
	case 1:
		removeSecret(secret)
	default:
		return
	}
}

func editSecret(s models.Secret) {
	unprotectName, err := crypto.UnCrypt(s.Data1)
	if err != nil {
		fmt.Println("Error when unprotect data, try again")
		return
	}

	_, err = crypto.UnCrypt(s.Data2)
	if err != nil {
		fmt.Println("Error when unprotect data, try again")
		return
	}

	uprotectMeta, err := crypto.UnCrypt(s.Data3)
	if err != nil {
		fmt.Println("Error when unprotect data, try again")
		return
	}

	cli.ChangeData(secretNameLabel, &unprotectName)

	switch s.Type {
	case passLoginType:
		/*tmp, err := utils.FromGOB64(unprotectName)
		if err != nil {

		}

		lp, ok := tmp.(LoginPass)
		if !ok {

		}*/

	default:

	}

	cli.ChangeData("", &uprotectMeta)
}

func removeSecret(s models.Secret) {
	err := network.RemoveSecret(s.Uid, s.Version)
	if err != nil {
		fmt.Println("Error when remove secret, try again")
	}
}

func printSecret(s models.Secret) error {
	unСryptData, err := crypto.UnCrypt(s.Data2)
	if err != nil {
		return err
	}

	fmt.Println(unСryptData)
	return nil
}
