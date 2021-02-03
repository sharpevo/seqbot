package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/sharpevo/seqbot/pkg/util"
)

const (
	CMD_WATCH = "watch"
	CMD_RUN   = "run"
	CMD_SEND  = "send"
)

type Mgi2000Command struct{}

func NewMgi2000Command() *Mgi2000Command {
	return &Mgi2000Command{}
}

func (m *Mgi2000Command) validate() error {
	if len(os.Args) < 2 {
		m.Usage()
		return fmt.Errorf("Error: command is required")
	}
	return nil
}

func (m *Mgi2000Command) Execute() error {
	if err := m.validate(); err != nil {
		return err
	}
	watchFlagSet := flag.NewFlagSet(CMD_WATCH, flag.ExitOnError)
	runFlagSet := flag.NewFlagSet(CMD_RUN, flag.ExitOnError)
	sendFlagSet := flag.NewFlagSet(CMD_SEND, flag.ExitOnError)
	switch os.Args[1] {
	case CMD_WATCH:
		watchCommand := NewWatchCommand(watchFlagSet)
		if err := watchFlagSet.Parse(os.Args[2:]); err != nil {
			watchFlagSet.PrintDefaults()
			return err
		}
		return watchCommand.Execute()
	case CMD_RUN:
		runCommand := NewRunCommand(runFlagSet)
		if err := runFlagSet.Parse(os.Args[2:]); err != nil {
			runFlagSet.PrintDefaults()
			return err
		}
		return runCommand.Execute()
	case CMD_SEND:
		sendCommand := util.NewSendCommand(sendFlagSet)
		if err := sendFlagSet.Parse(os.Args[2:]); err != nil {
			sendFlagSet.PrintDefaults()
			return err
		}
		return sendCommand.Execute()
	default:
		m.Usage()
		return fmt.Errorf("Error: invalid command")
	}
	return nil
}

func (t *Mgi2000Command) Usage() {
	fmt.Printf(`%s manages data generated from sequencer and sends notifications.

Commands:
  run: run actions.
  send: send messages.
  watch: watch WFQLog, take actions and send messages.

Use "%s <command> -h" for more information about a given command.

`,
		os.Args[0],
		os.Args[0],
	)
}
