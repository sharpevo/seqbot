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

type Options struct {
	Path       string
	DingTokens arrayFlag
}

func AttachOptions(cmd *flag.FlagSet) *Options {
	options := &Options{}
	cmd.StringVar(
		&options.Path,
		"dir",
		"",
		"wfqlog path to watch",
	)
	cmd.Var(
		&options.DingTokens,
		"dingtoken",
		"token of DingTalk robots",
	)
	return options
}
