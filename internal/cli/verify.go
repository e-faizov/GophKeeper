package cli

import (
	"errors"
	"fmt"
)

func NotEmpty(s string) error {
	if len(s) == 0 {
		return errors.New("not empty")
	}
	return nil
}

func MoreThan(n int) func(string) error {
	return func(s string) error {
		if len(s) < n {
			return fmt.Errorf("less than %d", n)
		}
		return nil
	}
}
