package options

import (
	"flag"
)

type WfqOptions struct {
	WfqLogPath string
}

func AttachWfqOptions(cmd *flag.FlagSet) *WfqOptions {
	options := &WfqOptions{}
	cmd.StringVar(
		&options.WfqLogPath,
		"wfqlog",
		"",
		"wfqlog path to watch",
	)
	return options
}
