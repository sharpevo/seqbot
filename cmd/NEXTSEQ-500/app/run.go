package app

import (
	"flag"
	"fmt"

	"github.com/sharpevo/seqbot/cmd/options"
	"github.com/sharpevo/seqbot/internal/pkg/action"
	"github.com/sharpevo/seqbot/internal/pkg/sequencer"

	"github.com/sirupsen/logrus"
)

type RunCommand struct {
	flagPath  string
	sequencer sequencer.SequencerInterface
	actions   []action.ActionInterface

	debugOption  *options.DebugOptions
	logOption    *options.LogOptions
	actionOption *options.ActionOptions
}

func NewRunCommand(flagSet *flag.FlagSet) *RunCommand {
	cmd := &RunCommand{
		sequencer: &sequencer.Nextseq500{},

		debugOption:  options.AttachDebugOptions(flagSet),
		logOption:    options.AttachLogOptions(flagSet),
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
	if err := r.logOption.Init(); err != nil {
		return err
	}
	if r.flagPath == "" {
		return fmt.Errorf("flagpath is required")
	}
	r.actions = []action.ActionInterface{
		&action.BarcodeAction{},
		&action.SlideAction{},
	}
	actions, err := r.actionOption.Actions()
	if err != nil {
		return err
	}
	for _, a := range actions {
		r.actions = append(r.actions, a)
		logrus.Infof("add action %s", a.Name())
	}
	if r.debugOption.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

func (r *RunCommand) Execute() error {
	if err := r.validate(); err != nil {
		return err
	}
	logrus.Info("manually running actions started")
	defer logrus.Info("manually running actions done")
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
