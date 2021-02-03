package options

import (
	"flag"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type LogOptions struct {
	LogPath string
}

func AttachLogOptions(cmd *flag.FlagSet) *LogOptions {
	options := &LogOptions{}
	cmd.StringVar(
		&options.LogPath,
		"log",
		"seqbot.log",
		"log file",
	)
	return options
}

func (l *LogOptions) Init() error {
	logFile, err := os.OpenFile(l.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:          true,
		DisableLevelTruncation: false,
	})
	return nil
}
