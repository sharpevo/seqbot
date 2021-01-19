package action

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sharpevo/seqbot/pkg/util"
)

const (
	LAYOUT_ARCHIVE      = "200601"
	NAME_ARCHIVE_ACTION = "Archive"

	MSG_TPL_ARCHIVE_SUCC = "- %s: %s\n"
	MSG_TPL_ARCHIVE_FAIL = "- %s: failed\n"
)

type ArchiveAction struct{}

func (a *ArchiveAction) Run(
	wfqLogPath string,
	chipId string,
) (string, error) {
	rootPath := util.ResultRootPathFromWFQLogPath(wfqLogPath)
	archivePath, err := getArchivePath(rootPath, time.Now())
	if err != nil {
		return fmt.Sprintf(MSG_TPL_ARCHIVE_FAIL, NAME_ARCHIVE_ACTION), err
	}
	return fmt.Sprintf(
			MSG_TPL_ARCHIVE_SUCC, NAME_ARCHIVE_ACTION, getOutput(archivePath),
		), os.Rename(
			filepath.Join(rootPath, chipId),
			filepath.Join(archivePath, chipId),
		)
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

func (a *ArchiveAction) Name() string {
	return NAME_ARCHIVE_ACTION
}
