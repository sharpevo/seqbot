package options

import (
	"flag"
	"fmt"
)

const (
	ADAPTER_INOTIFY = "inotify"
	ADAPTER_SCAN    = "scan"
)

type Nextseq500Options struct {
	DataPath     string
	Adapter      string
	ScanInterval uint
}

func AttachNextseq500Options(cmd *flag.FlagSet) *Nextseq500Options {
	options := &Nextseq500Options{}
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
