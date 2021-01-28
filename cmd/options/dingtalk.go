package options

import (
	"flag"
)

type DingtalkOptions struct {
	DingTokens arrayFlag
}

func AttachDingtalkOptions(cmd *flag.FlagSet) *DingtalkOptions {
	options := &DingtalkOptions{}
	cmd.Var(
		&options.DingTokens,
		"dingtoken",
		"token of DingTalk robots",
	)
	return options
}
