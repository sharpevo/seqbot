package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	nextOptions "github.com/sharpevo/seqbot/cmd/NEXTSEQ-500/app/options"
	"github.com/sharpevo/seqbot/cmd/options"
	"github.com/sharpevo/seqbot/internal/pkg/action"
	"github.com/sharpevo/seqbot/internal/pkg/sequencer"
	"github.com/sharpevo/seqbot/pkg/messenger"
	"github.com/sharpevo/seqbot/pkg/util"

	"github.com/sirupsen/logrus"
)

const (
	MSG_PLATFORM = "###### NEXTSEQ-500"
)

type WatchCommand struct {
	sequencer  sequencer.SequencerInterface
	messengers []messenger.Messenger
	actions    []action.ActionInterface

	option         *nextOptions.Nextseq500Options
	debugOption    *options.DebugOptions
	logOption      *options.LogOptions
	dingtalkOption *options.DingtalkOptions
	actionOption   *options.ActionOptions
}

func NewWatchCommand(flagSet *flag.FlagSet) *WatchCommand {
	return &WatchCommand{
		sequencer: &sequencer.Nextseq500{},

		option:         nextOptions.AttachNextseq500Options(flagSet),
		debugOption:    options.AttachDebugOptions(flagSet),
		logOption:      options.AttachLogOptions(flagSet),
		dingtalkOption: options.AttachDingtalkOptions(flagSet),
		actionOption:   options.AttachActionOptions(flagSet),
	}
}

func (w *WatchCommand) validate() error {
	if err := w.logOption.Init(); err != nil {
		return err
	}
	if w.option.DataPath == "" {
		return fmt.Errorf("data path is required")
	}
	for _, token := range w.dingtalkOption.DingTokens {
		dingbot := messenger.NewDingBot(token)
		w.messengers = append(w.messengers, dingbot)
		logrus.Infof("add messenger: %s", dingbot)
	}
	w.actions = []action.ActionInterface{
		&action.BarcodeAction{},
		&action.SlideAction{},
	}
	actions, err := w.actionOption.Actions()
	if err != nil {
		return err
	}
	for _, a := range actions {
		w.actions = append(w.actions, a)
		logrus.Infof("add action %s", a.Name())
	}
	if w.debugOption.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

func (w *WatchCommand) Execute() error {
	if err := w.validate(); err != nil {
		return err
	}
	defer logrus.Info("watching done")
	switch w.option.Adapter {
	case nextOptions.ADAPTER_INOTIFY:
		return fmt.Errorf("NextSeq 500 does not support inotify adapter")
	case nextOptions.ADAPTER_SCAN:
		logrus.Info("scanning started")
		return w.scan()
	default:
		return fmt.Errorf("invalid adapter")
	}
}

type seenMap map[string]struct{}

func (s seenMap) addFile(filePath string) bool {
	key := filepath.Join(filepath.Base(filepath.Dir(filePath)), filepath.Base(filePath))
	_, ok := s[key]
	if !ok {
		s[key] = struct{}{}
		logrus.Infof("seen %s", key)
	}
	return !ok
}

func (w *WatchCommand) scan() error {
	seen := seenMap{}
	err := filepath.Walk(
		w.option.DataPath,
		w.checkDir(seen, nil),
	)
	if err != nil {
		return err
	}
	for {
		err := filepath.Walk(
			w.option.DataPath,
			w.checkDir(seen, w.process),
		)
		if err != nil {
			logrus.Errorf(
				"failed to scan %s: %v", w.option.DataPath, err)
		}
		time.Sleep(time.Duration(w.option.ScanInterval) * time.Second)
	}
}

func (w *WatchCommand) process(filePath string) {
	slide, err := w.Sequencer().GetSlide(filePath)
	if err != nil {
		logrus.Errorf("failed to parse slide: %s", filePath)
		return
	}
	msg := util.NewMessage("\n")
	for _, a := range w.actions {
		output, err := a.Run(filePath, w)
		if err != nil {
			logrus.Errorf("failed to run '%s' on '%s': %v",
				a.Name(), slide, err)
		} else {
			logrus.Infof("action '%s' on '%s' success: %s",
				a.Name(), slide, output)
		}
		msg.Add(output)
	}
	w.send(msg.String())
	return
}

func (w *WatchCommand) send(message string) {
	if message == "" {
		logrus.Warn("empty message is ignored")
		return
	}
	message = fmt.Sprintf("%s\n%s", message, MSG_PLATFORM)
	for _, messenger := range w.messengers {
		err := messenger.Send(message)
		if err != nil {
			logrus.Errorf("failed to send message by %s: %v", messenger, err)
		}
	}
}

func (w *WatchCommand) Sequencer() sequencer.SequencerInterface {
	return w.sequencer
}

func (w *WatchCommand) checkDir(seen seenMap, process func(string)) filepath.WalkFunc {
	return func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsDir() {
			if util.IsArchiveDir(p) {
				logrus.Debugf("ignore archive directory: %s", p)
				return filepath.SkipDir
			}
			logrus.Debugf("ignore dir: %s", p)
			return nil
		}
		success, err := w.Sequencer().IsSuccess(p)
		if err != nil || !success {
			logrus.Debugf("ignore file: %s (%v)", p, err)
			return nil
		}
		if seen.addFile(p) && process != nil {
			process(p)
		}
		return nil
	}
}
