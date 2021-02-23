package app

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/sharpevo/seqbot/cmd/options"
	"github.com/sharpevo/seqbot/pkg/messenger"

	"github.com/sirupsen/logrus"
)

type SendCommand struct {
	msgFile string

	messengers     []messenger.Messenger
	logOption      *options.LogOptions
	dingtalkOption *options.DingtalkOptions
}

func NewSendCommand(flagSet *flag.FlagSet) *SendCommand {
	cmd := &SendCommand{
		logOption:      options.AttachLogOptions(flagSet),
		dingtalkOption: options.AttachDingtalkOptions(flagSet),
	}
	flagSet.StringVar(
		&cmd.msgFile,
		"msgfile",
		"",
		"message file to send",
	)
	return cmd
}

func (s *SendCommand) validate() error {
	if err := s.logOption.Init(); err != nil {
		return err
	}
	if s.msgFile == "" {
		return fmt.Errorf("message file is required")
	}
	if len(s.dingtalkOption.DingTokens) == 0 {
		return fmt.Errorf("token is required")
	}
	for _, token := range s.dingtalkOption.DingTokens {
		dingbot := messenger.NewDingBot(token)
		s.messengers = append(s.messengers, dingbot)
		logrus.Infof("add messenger: %s", dingbot)
	}
	return nil
}

func (s *SendCommand) Execute() error {
	if err := s.validate(); err != nil {
		return err
	}
	logrus.Info("manually sending notification started")
	defer logrus.Info("manually sending notification done")
	content, err := ioutil.ReadFile(s.msgFile)
	if err != nil {
		return err
	}
	for _, messenger := range s.messengers {
		err := messenger.Send(string(content))
		if err != nil {
			logrus.Errorf("failed to send message by %s: %v", messenger, err)
		}
		logrus.Infof("message sent: %s", string(content))
	}
	return nil
}
