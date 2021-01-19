package action

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sharpevo/seqbot/pkg/util"
)

const (
	NAME_ARCHIVE   = "Archive"
	LAYOUT_ARCHIVE = "200601"

	MSG_TPL_ARCHIVE_SUCC = "- Archive: %s"
	MSG_ARCHIVE_FAIL     = "- Archive: -"
)

type ArchiveAction struct{}

func (a *ArchiveAction) Run(
	eventName string,
	wfqLogPath string,
	chipId string,
) (string, error) {
	rootPath := util.ResultRootPathFromWFQLogPath(wfqLogPath)
	archivePath, err := getArchivePath(rootPath, time.Now())
	if err != nil {
		return MSG_ARCHIVE_FAIL, err
	}
	err = os.Rename(
		filepath.Join(rootPath, chipId),
		filepath.Join(archivePath, chipId))
	if err != nil {
		return MSG_ARCHIVE_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_ARCHIVE_SUCC, getOutput(archivePath)), nil
}

func (a *ArchiveAction) Name() string {
	return NAME_ARCHIVE
}

func getArchivePath(rootPath string, timestamp time.Time) (string, error) {
	archivePath := filepath.Join(rootPath, getArchiveName(timestamp))
	return archivePath, os.MkdirAll(archivePath, 0755)
}

func getArchiveName(timestamp time.Time) string {
	return timestamp.Format(LAYOUT_ARCHIVE)
}

func getOutput(archivePath string) string {
	return filepath.Base(archivePath)
}
