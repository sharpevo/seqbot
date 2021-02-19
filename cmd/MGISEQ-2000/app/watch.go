package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	logOption      *options.LogOptions
	dingtalkOption *options.DingtalkOptions
	actionOption   *options.ActionOptions
}

func NewWatchCommand(flagSet *flag.FlagSet) *WatchCommand {
	return &WatchCommand{
		sequencer: &sequencer.Mgiseq2000{},

		option:         mgiOptions.AttachMgiseq2000Options(flagSet),
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
	if err := w.validate(); err != nil {
		return err
	}
	defer logrus.Info("watching done")
	switch w.option.Adapter {
	case mgiOptions.ADAPTER_INOTIFY:
		logrus.Info("watching started")
		return w.watch()
	case mgiOptions.ADAPTER_SCAN:
		logrus.Info("scanning started")
		return w.scan()
	default:
		return fmt.Errorf("invalid adapter")
	}
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
						if util.IsArchiveDir(event.Name) {
							logrus.Infof("ignore archive directory: %s", event.Name)
							continue
						}
						watcher.Add(event.Name)
						logrus.Infof("watching directory: %s", event.Name)
						continue
					}
					w.process(event.Name)
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
				if util.IsArchiveDir(p) {
					logrus.Infof("ignore archive directory: %s", p)
					return nil
				}
				watcher.Add(p)
				logrus.Infof("watching directory: %s", p)
			}
			return nil
		})
	<-done
	return nil
}

type seenMap map[string]struct{}

func (s seenMap) addFile(filePath string) bool {
	key := filepath.Base(filePath)
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
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsDir() {
				return nil
			}
			seen.addFile(p)
			return nil
		},
	)
	if err != nil {
		return err
	}
	for {
		err := filepath.Walk(
			w.option.DataPath,
			func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.Mode().IsDir() {
					return nil
				}
				if seen.addFile(p) {
					w.process(p)
				}
				return nil
			})
		if err != nil {
			logrus.Errorf(
				"failed to scan %s: %v", w.option.DataPath, err)
		}
		time.Sleep(time.Duration(w.option.ScanInterval) * time.Second)
	}
}

func (w *WatchCommand) process(filePath string) {
	success, err := w.Sequencer().IsSuccess(filePath)
	if err != nil {
		logrus.Errorf("failed to check success %s: %v", filePath, err)
		return
	}
	if !success {
		logrus.Debugf("ignore event: %s", filePath)
		return
	}
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
