package options

import (
	"flag"
)

type DebugOptions struct {
	Debug bool
}

func AttachDebugOptions(cmd *flag.FlagSet) *DebugOptions {
	options := &DebugOptions{}
	cmd.BoolVar(
		&options.Debug,
		"debug",
		false,
		"show debug message",
	)
	return options
}
