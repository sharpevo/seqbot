package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sharpevo/seqbot/cmd/DNBSEQ-T7/app/options"
	"github.com/sharpevo/seqbot/internal/pkg/action"
	"github.com/sharpevo/seqbot/internal/pkg/lane"
	"github.com/sharpevo/seqbot/pkg/messenger"
	"github.com/sharpevo/seqbot/pkg/util"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

const (
	DIR_RUNNING = "running"
	DIR_FINISH  = "finish"
	DIR_FAIL    = "fail"

	MSG_PLATFORM = "###### DNBSEQ-T7"
)

type WatchCommand struct {
	messengers []messenger.Messenger
	actions    []action.ActionInterface

	option         *options.Options
	wfqOption      *options.WfqOptions
	dingtalkOption *options.DingtalkOptions
	actionOption   *options.ActionOptions
}

func NewWatchCommand(flagSet *flag.FlagSet) *WatchCommand {
	return &WatchCommand{
		option:         options.AttachOptions(flagSet),
		wfqOption:      options.AttachWfqOptions(flagSet),
		dingtalkOption: options.AttachDingtalkOptions(flagSet),
		actionOption:   options.AttachActionOptions(flagSet),
	}
}

func (w *WatchCommand) validate() error {
	flag.Parse()
	if w.wfqOption.WfqLogPath == "" {
		return fmt.Errorf("wfqlog is required")
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
	if w.actionOption.ActionSummary {
		summaryAction := &action.SummaryAction{}
		w.actions = append(w.actions, summaryAction)
		logrus.Infof("add action %s", summaryAction.Name())
	}
	if w.actionOption.ActionWfqTime {
		wfqTimeAction := &action.WfqTimeAction{}
		w.actions = append(w.actions, wfqTimeAction)
		logrus.Infof("add action %s", wfqTimeAction.Name())
	}
	if w.actionOption.ActionArchive {
		archiveAction := &action.ArchiveAction{}
		w.actions = append(w.actions, archiveAction)
		logrus.Infof("add action %s", archiveAction.Name())
	}
	if w.option.Debug {
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
		return fmt.Errorf("failed to add watcher '%s': %s", w.wfqOption.WfqLogPath, err)
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
					chipId := util.ChipIdFromFlagPath(event.Name)
					message, err := w.update(event.Name, chipId)
					if err != nil {
						logrus.Errorf("failed to update %s: %v", chipId, err)
						message = fmt.Sprintf(
							"**%s**: sequencing completed, but failed to update database.",
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

	dirFinish := filepath.Join(w.wfqOption.WfqLogPath, DIR_FINISH)
	dirRunning := filepath.Join(w.wfqOption.WfqLogPath, DIR_RUNNING)
	dirFail := filepath.Join(w.wfqOption.WfqLogPath, DIR_FAIL)
	watcher.Add(dirFinish)
	watcher.Add(dirRunning)
	watcher.Add(dirFail)
	logrus.Infof("watching directories: %s, %s, %s", dirFinish, dirRunning, dirFail)
	<-done
	return nil
}

func (w *WatchCommand) update(eventName string, chipId string) (string, error) {
	msg := &Message{sep: "\n"}
	dir := filepath.Base(filepath.Dir(eventName))
	switch dir {
	case DIR_RUNNING:
		l := lane.NewLane(chipId)
		if err := l.Start(); err != nil {
			return msg.String(), err
		}
		logrus.Infof("%s: WFQ has been started.", l.ChipId)
		return "", nil
	case DIR_FINISH:
		for _, a := range w.actions {
			output, err := a.Run(eventName, w.wfqOption.WfqLogPath, chipId)
			if err != nil {
				logrus.Errorf("failed to run '%s' on '%s': %v", a.Name(), chipId, err)
			} else {
				logrus.Infof("action '%s' on '%s' success: %s", a.Name(), chipId, output)
			}
			msg.Add(output)
		}
	case DIR_FAIL:
		if w.isExistRunningOrDuplicateLane(eventName) {
			logrus.Warnf("Not mark as fail for exist running or duplicate lane")
			return fmt.Sprintf("**%s**: Job exists or duplicate lane.", chipId), nil
		}
		l := lane.NewLane(chipId)
		if err := l.Finish(); err != nil {
			return msg.String(), err
		}
		return fmt.Sprintf("**%s**: WFQ failed, %s.", l.ChipId, l.Duration()), nil
	default:
		return msg.String(), fmt.Errorf("invalide dir: %s", dir)
	}
	return msg.String(), nil
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

func (w *WatchCommand) isExistRunningOrDuplicateLane(failedFlagPath string) bool {
	_, err := os.Stat(filepath.Join(
		w.wfqOption.WfqLogPath,
		DIR_RUNNING,
		filepath.Base(failedFlagPath)))
	return !os.IsNotExist(err)
}
