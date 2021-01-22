package options

import (
	"flag"
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

type DingtalkOptions struct {
	DingTokens arrayFlag
}

func AttachDingtalkOptions(cmd *flag.FlagSet) *DingtalkOptions {
	options := &DingtalkOptions{}
	cmd.Var(
		&options.DingTokens,
		"dingtoken",
		"token of DingTalk robots",
	)
	return options
}
