package options

import (
	"flag"
	"fmt"

	"github.com/sharpevo/seqbot/internal/pkg/action"
)

const (
	OPTION_ACTION_ARCHIVE    = "archive"
	OPTION_ACTION_WFQTIME    = "wfqtime"
	OPTION_ACTION_UPLOADTIME = "uploadtime"
	OPTION_ACTION_EXTRA      = "misc"
)

type ActionOptions struct {
	actions arrayFlag
}

func AttachActionOptions(cmd *flag.FlagSet) *ActionOptions {
	options := &ActionOptions{}
	cmd.Var(
		&options.actions,
		"action",
		"actions when event captured",
	)
	return options
}

func (o *ActionOptions) Actions() ([]action.ActionInterface, error) {
	actions := []action.ActionInterface{}
	for _, a := range o.actions {
		switch a {
		case OPTION_ACTION_ARCHIVE:
			actions = append(actions, &action.ArchiveAction{})
		case OPTION_ACTION_WFQTIME:
			actions = append(actions, &action.WfqTimeAction{})
		case OPTION_ACTION_UPLOADTIME:
			actions = append(actions, &action.UploadTimeAction{})
		case OPTION_ACTION_EXTRA:
			actions = append(actions, &action.ExtraAction{})
		default:
			return actions, fmt.Errorf("invalid action '%s'", a)
		}
	}
	return actions, nil
}
