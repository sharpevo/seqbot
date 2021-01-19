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
	WfqLogPath    string
	DingTokens    arrayFlag
	Debug         bool
	ActionArchive bool
	ActionSummary bool
	ActionWfqTime bool
}

func AttachOptions(cmd *flag.FlagSet) *Options {
	options := &Options{}
	cmd.StringVar(
		&options.WfqLogPath,
		"wfqlog",
		"",
		"wfqlog path to watch",
	)
	cmd.BoolVar(
		&options.Debug,
		"debug",
		false,
		"show debug message",
	)
	cmd.BoolVar(
		&options.ActionArchive,
		"actionarchive",
		true,
		"archive result in YYYYMM directory",
	)
	cmd.BoolVar(
		&options.ActionSummary,
		"actionsummary",
		true,
		"summary the number and size of fq.gz",
	)
	cmd.BoolVar(
		&options.ActionWfqTime,
		"actionwfqtime",
		true,
		"calculate the time of writing fastq",
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
			"-wfqlog=/path/to/wfqlog -dingtoken=token1 -dingtolken=token2",
		)
		flag.PrintDefaults()
	}
	return options
}
