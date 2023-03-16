package cli

import "fmt"

type SelectionItem string

func Selection(items []SelectionItem) int {
	var selection int
	for {
		for i, v := range items {
			fmt.Println(i+1, v)
		}
		_, err := fmt.Scan(&selection)
		if err != nil {
			fmt.Println(err, "try again")
			continue
		}
		if selection > len(items) {
			fmt.Println("wrong number try again")
			continue
		}
		break
	}
	return selection - 1
}

type EnterDataItem struct {
	Name   string
	Data   *string
	Verify func(string) error
}

func EnterData(data []EnterDataItem) {
	for _, v := range data {
		if v.Data == nil {
			panic("Undefined behavior")
		}
		for {
			fmt.Print(v.Name + ": ")
			var tmp string
			fmt.Scan(&tmp)
			if v.Verify != nil {
				err := v.Verify(tmp)
				if err != nil {
					fmt.Println(err, "try again")
					continue
				}
			}
			*v.Data = tmp
			break
		}

	}
}

func ChangeData(name string, data *string) {
	if data == nil {
		return
	}

	if len(*data) != 0 {
		fmt.Println("Current value:", *data)
	}
	fmt.Println("Change? yes(y):no(n)")
	var choice string
	fmt.Scan(&choice)

	if choice == "n" || choice == "no" {
		return
	}

	fields := []EnterDataItem{
		{
			Name:   name,
			Data:   data,
			Verify: NotEmpty,
		},
	}

	EnterData(fields)
}
