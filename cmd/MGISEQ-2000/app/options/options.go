package options

import (
	"flag"
	"fmt"
	"os"
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
	DataPath   string
	DingTokens arrayFlag
	Debug      bool
}

func AttachOptions(cmd *flag.FlagSet) *Options {
	options := &Options{}
	cmd.StringVar(
		&options.DataPath,
		"data",
		"",
		"data path to watch",
	)
	cmd.BoolVar(
		&options.Debug,
		"debug",
		false,
		"show debug message",
	)
	cmd.Var(
		&options.DingTokens,
		"dingtoken",
		"token of DingTalk robots",
	)
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"Usage of %s:\n  %s %s\n",
			os.Args[0],
			os.Args[0],
			"-data=/path/to/data -dingtoken=token1 -dingtolken=token2",
		)
		flag.PrintDefaults()
	}
	return options
}
