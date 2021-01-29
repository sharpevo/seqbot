package options

import (
	"flag"
)

type Dnbseqt7Options struct {
	WfqLogPath string
}

func AttachDnbseqt7Options(cmd *flag.FlagSet) *Dnbseqt7Options {
	options := &Dnbseqt7Options{}
	cmd.StringVar(
		&options.WfqLogPath,
		"wfqlog",
		"",
		"wfqlog path to watch",
	)
	return options
}
