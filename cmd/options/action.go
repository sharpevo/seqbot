package options

import (
	"flag"
)

type ActionOptions struct {
	ActionArchive bool
	ActionSummary bool
	ActionWfqTime bool
}

func AttachActionOptions(cmd *flag.FlagSet) *ActionOptions {
	options := &ActionOptions{}
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
	return options
}
