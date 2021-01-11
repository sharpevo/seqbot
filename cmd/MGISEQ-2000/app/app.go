package app

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
							"**%s** sequencing completed.\n- Slide: %s\n",
							dnbId,
							slideId)
						m.send(message)
						logrus.Infof("message sent: %s", message)
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

	watcher.Add(m.dataPath)
	logrus.Infof("watching directory: %s", m.dataPath)
	files, err := ioutil.ReadDir(m.dataPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		filePath := filepath.Join(m.dataPath, f.Name())
		if f.Mode().IsDir() {
			watcher.Add(filePath)
			logrus.Infof("watching directory: %s", filePath)
			continue
		}
	}
	<-done
	return nil
}

func (m *Mgi2000Command) send(message string) {
	for _, messenger := range m.messengers {
		err := messenger.Send(message)
		if err != nil {
			logrus.Errorf("failed to send message by %s: %v", messenger, err)
		}
	}
}
