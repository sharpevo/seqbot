package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sharpevo/seqbot/cmd/MGISEQ-2000/app/options"
	"github.com/sharpevo/seqbot/pkg/messenger"
	"github.com/sharpevo/seqbot/pkg/util"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

const (
	FILE_LOG = "seqbot.log"

	MSG_PLATFORM = "###### MGISEQ-2000"
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

type Mgi2000Command struct {
	dataPath   string
	options    *options.Options
	messengers []messenger.Messenger
}

func NewMgi2000Command() *Mgi2000Command {
	return &Mgi2000Command{
		options: options.AttachOptions(flag.CommandLine),
	}
}

func (m *Mgi2000Command) validate() error {
	flag.Parse()
	if m.options.DataPath == "" {
		return fmt.Errorf("data path is required")
	}
	for _, token := range m.options.DingTokens {
		dingbot := messenger.NewDingBot(token)
		m.messengers = append(m.messengers, dingbot)
		logrus.Infof("add messenger: %s", dingbot)
	}
	m.dataPath = m.options.DataPath
	if m.options.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

func (m *Mgi2000Command) Execute() error {
	if err := m.validate(); err != nil {
		return err
	}
	return m.watch()
}

func (m *Mgi2000Command) watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to add watcher '%s': %s", m.dataPath, err)
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

					if util.IsSuccess(event.Name) {
						slideId, dnbId := util.ParseMgiInfo(event.Name)
						message := fmt.Sprintf(
							"**%s**: sequencing completed.\n- Slide: %s\n",
							dnbId,
							slideId)
						m.send(message)
						logrus.Infof("message sent: %s", message)
					} else {
						logrus.Infof("ignore event: %s", event.Name)
					}
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
		m.dataPath,
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

func (m *Mgi2000Command) send(message string) {
	message = fmt.Sprintf("%s\n%s", message, MSG_PLATFORM)
	for _, messenger := range m.messengers {
		err := messenger.Send(message)
		if err != nil {
			logrus.Errorf("failed to send message by %s: %v", messenger, err)
		}
	}
}
