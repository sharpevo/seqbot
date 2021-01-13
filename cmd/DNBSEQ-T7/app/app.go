package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sharpevo/seqbot/cmd/DNBSEQ-T7/app/options"
	"github.com/sharpevo/seqbot/internal/pkg/flagjson"
	"github.com/sharpevo/seqbot/internal/pkg/lane"
	"github.com/sharpevo/seqbot/pkg/messenger"
	"github.com/sharpevo/seqbot/pkg/util"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

const (
	FILE_LOG = "seqbot.log"

	DIR_RUNNING = "running"
	DIR_FINISH  = "finish"
	DIR_FAIL    = "fail"

	MSG_PLATFORM = "###### DNBSEQ-T7"
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
	if w.options.WfqLogPath == "" {
		return fmt.Errorf("wfqlog is required")
	}
	for _, token := range w.options.DingTokens {
		dingbot := messenger.NewDingBot(token)
		w.messengers = append(w.messengers, dingbot)
		logrus.Infof("add messenger: %s", dingbot)
	}
	w.watchDir = w.options.WfqLogPath
	if w.options.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
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
					if filepath.Ext(event.Name) != ".json" ||
						strings.HasPrefix(filepath.Base(event.Name), ".") {
						logrus.Debugf("ignore creation: %s", event.Name)
						continue
					}
					chipId := getChipId(event.Name)
					message, err := w.update(event.Name, chipId)
					if err != nil {
						logrus.Errorf("failed to update %s: %v", chipId, err)
						message = fmt.Sprintf(
							"**%s**: WFQ completed, but failed to update database.",
							chipId)
						w.send(message)
						continue
					}
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

func (w *WatchCommand) update(eventName string, chipId string) (string, error) {
	message := ""
	dir := filepath.Base(filepath.Dir(eventName))
	switch dir {
	case DIR_RUNNING:
		l := lane.NewLane(chipId)
		if err := l.Start(); err != nil {
			return message, err
		}
		logrus.Infof("%s: WFQ has been started.", l.ChipId)
		return "", nil
	case DIR_FINISH:
		l := lane.NewLane(chipId)
		if err := l.Finish(); err != nil {
			return message, err
		}
		count, size, err := util.FastqCountAndSize(w.options.WfqLogPath, chipId)
		if err != nil {
			return message, err
		}
		f, err := flagjson.ReadFlag(eventName)
		if err != nil {
			return message, err
		}
		message = fmt.Sprintf(
			"**%s** sequencing completed, with %d fq.gz in %s.\n- Slide: %s\n- WFQ Time: %s",
			f.BarcodeType(), count, size, l.ChipId, l.Duration())
		if !w.options.Archive {
			logrus.Infof("ignore archiving %s", chipId)
			return message, nil
		}
		archivedPath, err := w.archive(chipId)
		if err != nil {
			logrus.Errorf("failed to archive %s: %v", chipId, err)
			message = fmt.Sprintf("%s\n- Archive: failed", message)
			return message, nil
		}
		logrus.Infof("archived %s: %s", chipId, archivedPath)
		message = fmt.Sprintf("%s\n- Archive: %s",
			message, filepath.Base(archivedPath))
		return message, nil
	case DIR_FAIL:
		if w.isExistRunningOrDuplicateLane(eventName) {
			logrus.Warnf("Not mark as fail for exist running or duplicate lane")
			return fmt.Sprintf("**%s**: Job exists or duplicate lane.", chipId), nil
		}
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
		logrus.Infof("message sent: %s", message)
	}
}

func (w *WatchCommand) archive(chipId string) (string, error) {
	resultPath := util.ResultRootPathFromWFQLogPath(w.options.WfqLogPath)
	archivedPath, err := util.CreateArchivedDir(resultPath, time.Now())
	if err != nil {
		return archivedPath, err
	}
	oldPath := filepath.Join(resultPath, chipId)
	if err := os.Rename(oldPath, filepath.Join(archivedPath, chipId)); err != nil {
		return archivedPath, err
	}
	return archivedPath, nil
}

func (w *WatchCommand) isExistRunningOrDuplicateLane(failedFlagPath string) bool {
	_, err := os.Stat(filepath.Join(
		w.options.WfqLogPath,
		DIR_RUNNING,
		filepath.Base(failedFlagPath)))
	return !os.IsNotExist(err)
}

func getChipId(filename string) string {
	return strings.Split(filepath.Base(filename), "_")[0]
}
