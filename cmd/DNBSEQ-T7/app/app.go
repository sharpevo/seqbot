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

type T7Command struct{}

func NewT7Command() *T7Command {
	return &T7Command{}
}

func (t *T7Command) validate() error {
	if len(os.Args) < 2 {
		t.Usage()
		return fmt.Errorf("Error: command is required")
	}
	return nil
}

func (t *T7Command) Usage() {
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

func (t *T7Command) Execute() error {
	if err := t.validate(); err != nil {
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
		t.Usage()
		return fmt.Errorf("Error: invalid command")
	}
	return nil
}
