package options

import (
	"strings"
)

type arrayFlag []string

func (a *arrayFlag) String() string {
	return strings.Join(*a, ",")
}

func (a *arrayFlag) Set(value string) error {
	*a = append(*a, value)
	return nil
}
