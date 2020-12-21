package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sharpevo/seqbot/cmd/watch/app/options"
	"github.com/sharpevo/seqbot/internal/pkg/lane"
	"github.com/sharpevo/seqbot/pkg/messenger"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

const (
	FILE_LOG = "seqbot.log"

	DIR_RUNNING = "running"
	DIR_FINISH  = "finish"
	DIR_FAIL    = "fail"
)

func init() {
	logFile, err := os.OpenFile(FILE_LOG, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("failed to open log file", FILE_LOG)
		return
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:          true,
		DisableLevelTruncation: false,
	})
}

type WatchCommand struct {
	dingUrl    string
	watchDir   string
	options    *options.Options
	messengers []messenger.Messenger
}

func NewWatchCommand() *WatchCommand {
	return &WatchCommand{
		options: options.AttachOptions(flag.CommandLine),
	}
}

func (w *WatchCommand) validate() error {
	flag.Parse()
	for _, token := range w.options.DingTokens {
		dingbot := messenger.NewDingBot(token)
		w.messengers = append(w.messengers, dingbot)
		logrus.Infof("add messenger: %s", dingbot)
	}
	w.watchDir = w.options.Path
	return nil
}

func (w *WatchCommand) Execute() error {
	if err := w.validate(); err != nil {
		return err
	}
	return w.watch()
}

func (w *WatchCommand) watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to add watcher '%s': %s", w.watchDir, err)
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
					message, err := w.update(
						filepath.Base(filepath.Dir(event.Name)),
						getChipId(event.Name))
					if err != nil {
						logrus.Errorf("failed to update: %s", err)
					}
					logrus.Debugf("message composed: %s", message)
					w.send(message)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Errorf("watcher error: %s", err)
			}
		}
	}()

	dirFinish := filepath.Join(w.watchDir, DIR_FINISH)
	dirRunning := filepath.Join(w.watchDir, DIR_RUNNING)
	dirFail := filepath.Join(w.watchDir, DIR_FAIL)
	watcher.Add(dirFinish)
	watcher.Add(dirRunning)
	watcher.Add(dirFail)
	logrus.Infof("watching directories: %s, %s, %s", dirFinish, dirRunning, dirFail)
	<-done
	return nil
}

func (w *WatchCommand) update(dir string, chipId string) (string, error) {
	message := ""
	switch dir {
	case DIR_RUNNING:
		l := lane.NewLane(chipId)
		if err := l.Start(); err != nil {
			return message, err
		}
		return fmt.Sprintf("**%s**: WFQ has been started.", l.ChipId), nil
	case DIR_FINISH:
		l := lane.NewLane(chipId)
		if err := l.Finish(); err != nil {
			return message, err
		}
		return fmt.Sprintf("**%s**: WFQ completed, %s.", l.ChipId, l.Duration()), nil
	case DIR_FAIL:
		l := lane.NewLane(chipId)
		if err := l.Finish(); err != nil {
			return message, err
		}
		return fmt.Sprintf("**%s**: WFQ failed, %s.", l.ChipId, l.Duration()), nil
	default:
		return message, fmt.Errorf("invalide dir: %s", dir)
	}
}

func (w *WatchCommand) send(message string) {
	for _, messenger := range w.messengers {
		err := messenger.Send(message)
		if err != nil {
			logrus.Errorf("failed to send message by %s: %v", messenger, err)
		}
	}
}

func getChipId(filename string) string {
	return strings.Split(filepath.Base(filename), "_")[0]
}
