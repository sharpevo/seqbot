package options

import (
	"flag"
	"fmt"
)

const (
	ADAPTER_INOTIFY = "inotify"
	ADAPTER_SCAN    = "scan"
)

type Mgiseq2000Options struct {
	DataPath     string
	Adapter      string
	ScanInterval uint
}

func AttachMgiseq2000Options(cmd *flag.FlagSet) *Mgiseq2000Options {
	options := &Mgiseq2000Options{}
	cmd.StringVar(
		&options.DataPath,
		"data",
		"",
		"data path to watch",
	)
	cmd.StringVar(
		&options.Adapter,
		"adapter",
		ADAPTER_INOTIFY,
		fmt.Sprintf("%s or %s", ADAPTER_INOTIFY, ADAPTER_SCAN),
	)
	cmd.UintVar(
		&options.ScanInterval,
		"interval",
		300,
		"seconds to execute scanning",
	)
	return options
}
