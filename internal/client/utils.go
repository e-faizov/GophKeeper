package client

import (
	"fmt"

	"github.com/ShiraazMoollatjie/goluhn"

	"github.com/e-faizov/GophKeeper/internal/cli"
)

func dataMetaFields(tp int) (string, string) {
	var err error
	var data string
	for {
		switch tp {
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
		break
	}

	var secretMeta string
	metaField := []cli.EnterDataItem{
		{
			Name: metaLabel,
			Data: &secretMeta,
		},
	}
	cli.EnterData(metaField)
	return data, secretMeta
}
