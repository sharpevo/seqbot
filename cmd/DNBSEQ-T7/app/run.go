package app

import (
	"flag"
	"fmt"

	dnbseqt7Options "github.com/sharpevo/seqbot/cmd/DNBSEQ-T7/app/options"
	"github.com/sharpevo/seqbot/cmd/options"
	"github.com/sharpevo/seqbot/internal/pkg/action"
	"github.com/sharpevo/seqbot/internal/pkg/sequencer"

	"github.com/sirupsen/logrus"
)

type RunCommand struct {
	flagPath  string
	actions   []action.ActionInterface
	sequencer sequencer.SequencerInterface

	option       *dnbseqt7Options.Dnbseqt7Options
	debugOption  *options.DebugOptions
	actionOption *options.ActionOptions
}

func NewRunCommand(flagSet *flag.FlagSet) *RunCommand {
	cmd := &RunCommand{
		sequencer: &sequencer.Dnbseqt7{},

		option:       dnbseqt7Options.AttachDnbseqt7Options(flagSet),
		debugOption:  options.AttachDebugOptions(flagSet),
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
	slide, err := r.Sequencer().GetSlide(r.flagPath)
	if err != nil {
		logrus.Errorf("failed to parse slide: %s", r.flagPath)
		return err
	}
	for _, a := range r.actions {
		output, err := a.Run(r.flagPath, r)
		if err != nil {
			logrus.Errorf("failed to run '%s' on '%s': %v", a.Name(), slide, err)
		} else {
			logrus.Infof("action '%s' on '%s' success: %s", a.Name(), slide, output)
		}
	}
	return nil
}

func (r *RunCommand) Sequencer() sequencer.SequencerInterface {
	return r.sequencer
}
