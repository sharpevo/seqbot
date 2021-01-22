package options

import (
	"flag"
)

type Options struct {
	Debug bool
}

func AttachOptions(cmd *flag.FlagSet) *Options {
	options := &Options{}
	cmd.BoolVar(
		&options.Debug,
		"debug",
		false,
		"show debug message",
	)
	return options
}
