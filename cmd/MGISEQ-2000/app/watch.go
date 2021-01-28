package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	mgiOptions "github.com/sharpevo/seqbot/cmd/MGISEQ-2000/app/options"
	"github.com/sharpevo/seqbot/cmd/options"
	"github.com/sharpevo/seqbot/internal/pkg/action"
	"github.com/sharpevo/seqbot/internal/pkg/sequencer"
	"github.com/sharpevo/seqbot/pkg/messenger"
	"github.com/sharpevo/seqbot/pkg/util"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

const (
	MSG_PLATFORM = "###### MGISEQ-2000"
)

type WatchCommand struct {
	sequencer  sequencer.SequencerInterface
	messengers []messenger.Messenger
	actions    []action.ActionInterface

	option         *mgiOptions.Mgiseq2000Options
	debugOption    *options.DebugOptions
	dingtalkOption *options.DingtalkOptions
	actionOption   *options.ActionOptions
}

func NewWatchCommand(flagSet *flag.FlagSet) *WatchCommand {
	return &WatchCommand{
		sequencer: &sequencer.Mgiseq2000{},

		option:         mgiOptions.AttachMgiseq2000Options(flagSet),
		debugOption:    options.AttachDebugOptions(flagSet),
		dingtalkOption: options.AttachDingtalkOptions(flagSet),
		actionOption:   options.AttachActionOptions(flagSet),
	}
}

func (w *WatchCommand) validate() error {
	flag.Parse()
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
	if w.actionOption.ActionArchive {
		archiveAction := &action.ArchiveAction{}
		w.actions = append(w.actions, archiveAction)
		logrus.Infof("add action %s", archiveAction.Name())
	}
	if w.debugOption.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil

}

func (w *WatchCommand) Execute() error {
	logrus.Info("watching started")
	defer logrus.Info("watching done")
	if err := w.validate(); err != nil {
		return err
	}
	return w.watch()
}

func (w *WatchCommand) watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf(
			"failed to add watcher '%s': %v", w.option.DataPath, err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					f, _ := os.Stat(event.Name)
					if f.Mode().IsDir() {
						watcher.Add(event.Name)
						logrus.Infof("watching directory: %s", event.Name)
						continue
					}
					success, err := w.Sequencer().IsSuccess(event.Name)
					if err != nil {
						logrus.Errorf(
							"failed to check success %s: %v", event.Name, err)
						continue
					}
					if !success {
						logrus.Infof("ignore event: %s", event.Name)
						continue
					}
					slide, err := w.Sequencer().GetSlide(event.Name)
					if err != nil {
						logrus.Errorf("failed to parse slide: %s", event.Name)
						continue
					}
					msg := util.NewMessage("\n")
					for _, a := range w.actions {
						output, err := a.Run(event.Name, w)
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
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Errorf("watcher error: %s", err)
			}
		}
	}()

	err = filepath.Walk(
		w.option.DataPath,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsDir() {
				watcher.Add(p)
				logrus.Infof("watching directory: %s", p)
			}
			return nil
		})
	<-done
	return nil
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
