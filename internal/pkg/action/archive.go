package action

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	NAME_ARCHIVE = "Archive"

	MSG_TPL_ARCHIVE_SUCC = "- Archive: %s"
	MSG_ARCHIVE_FAIL     = "- Archive: -"
)

type ArchiveAction struct{}

func (a *ArchiveAction) Run(
	eventName string,
	command CommandInterface,
) (string, error) {
	srcDir, err := command.Sequencer().GetResultDir(eventName)
	if err != nil {
		return MSG_ARCHIVE_FAIL, err
	}
	dstDir, err := command.Sequencer().GetArchiveDir(eventName)
	if err != nil {
		return MSG_ARCHIVE_FAIL, err
	}
	slide, err := command.Sequencer().GetSlide(eventName)
	if err != nil {
		return MSG_ARCHIVE_FAIL, err
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return MSG_ARCHIVE_FAIL, err
	}
	if err := os.Rename(srcDir, filepath.Join(dstDir, slide)); err != nil {
		return MSG_ARCHIVE_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_ARCHIVE_SUCC, filepath.Base(dstDir)), nil
}

func (a *ArchiveAction) Name() string {
	return NAME_ARCHIVE
}
