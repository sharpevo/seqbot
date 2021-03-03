package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/sharpevo/seqbot/cmd/send/app"
)

const (
	CMD_WATCH = "watch"
	CMD_RUN   = "run"
	CMD_SEND  = "send"
)

type Nextseq500Command struct{}

func NewNextseq500Command() *Nextseq500Command {
	return &Nextseq500Command{}
}

func (n *Nextseq500Command) validate() error {
	if len(os.Args) < 2 {
		n.Usage()
		return fmt.Errorf("Error: command is required")
	}
	return nil
}

func (n *Nextseq500Command) Execute() error {
	if err := n.validate(); err != nil {
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
		sendCommand := app.NewSendCommand(sendFlagSet)
		if err := sendFlagSet.Parse(os.Args[2:]); err != nil {
			sendFlagSet.PrintDefaults()
			return err
		}
		return sendCommand.Execute()
	default:
		n.Usage()
		return fmt.Errorf("Error: invalid command")
	}
	return nil
}

func (n *Nextseq500Command) Usage() {
	fmt.Printf(`%s manages data generated from sequencer and sends notifications.

Commands:
  run: run actions.
  send: send messages.
  watch: watch flag file, take actions and send messages.

Use "%s <command> -h" for more information about a given command.

`,
		os.Args[0],
		os.Args[0],
	)
}
