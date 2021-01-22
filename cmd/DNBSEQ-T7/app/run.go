package app

import (
	"flag"
	"fmt"

	"github.com/sharpevo/seqbot/cmd/DNBSEQ-T7/app/options"
	"github.com/sharpevo/seqbot/internal/pkg/action"
	"github.com/sharpevo/seqbot/pkg/util"

	"github.com/sirupsen/logrus"
)

type RunCommand struct {
	option       *options.Options
	wfqOption    *options.WfqOptions
	actionOption *options.ActionOptions

	flagPath string
	actions  []action.ActionInterface
}

func NewRunCommand(flagSet *flag.FlagSet) *RunCommand {
	cmd := &RunCommand{
		option:       options.AttachOptions(flagSet),
		wfqOption:    options.AttachWfqOptions(flagSet),
		actionOption: options.AttachActionOptions(flagSet),
	}
	flagSet.StringVar(
		&cmd.flagPath,
		"flagpath",
		"",
		"flag file path",
	)
	return cmd
}

func (r *RunCommand) validate() error {
	if r.flagPath == "" {
		return fmt.Errorf("flagpath is required")
	}
	r.actions = []action.ActionInterface{
		&action.BarcodeAction{},
		&action.SlideAction{},
	}
	if r.actionOption.ActionSummary {
		summaryAction := &action.SummaryAction{}
		r.actions = append(r.actions, summaryAction)
		logrus.Infof("add action %s", summaryAction.Name())
	}
	if r.actionOption.ActionWfqTime {
		wfqTimeAction := &action.WfqTimeAction{}
		r.actions = append(r.actions, wfqTimeAction)
		logrus.Infof("add action %s", wfqTimeAction.Name())
	}
	if r.actionOption.ActionArchive {
		archiveAction := &action.ArchiveAction{}
		r.actions = append(r.actions, archiveAction)
		logrus.Infof("add action %s", archiveAction.Name())
	}
	return nil
}

func (r *RunCommand) Execute() error {
	logrus.Info("manually running actions started")
	defer logrus.Info("manually running actions done")
	if err := r.validate(); err != nil {
		return err
	}
	chipId := util.ChipIdFromFlagPath(r.flagPath)
	for _, a := range r.actions {
		output, err := a.Run(r.flagPath, r.wfqOption.WfqLogPath, chipId)
		if err != nil {
			logrus.Errorf("failed to run '%s' on '%s': %v", a.Name(), chipId, err)
		} else {
			logrus.Infof("action '%s' on '%s' success: %s", a.Name(), chipId, output)
		}
	}
	return nil
}
