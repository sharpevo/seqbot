package options

import (
	"flag"
)

type Mgiseq2000Options struct {
	DataPath string
}

func AttachMgiseq2000Options(cmd *flag.FlagSet) *Mgiseq2000Options {
	options := &Mgiseq2000Options{}
	cmd.StringVar(
		&options.DataPath,
		"data",
		"",
		"data path to watch",
	)
	return options
}
