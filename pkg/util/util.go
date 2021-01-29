package util

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	PATH_RESULT = "result/OutputFq/upload"
	PATH_FLAG   = "flag"

	LAYOUT_ARCHIVE = "200601"
)

func ResultRootPathFromWFQLogPath(wfqlogPath string) string {
	return filepath.Join(filepath.Dir(wfqlogPath), PATH_RESULT)
}

func ResultChipPathFromWFQLogPath(wfqlogPath string, chipId string) string {
	return filepath.Join(ResultRootPathFromWFQLogPath(wfqlogPath), chipId)
}

func FlagPathFromWFQLogPath(wfqlogPath string) string {
	return filepath.Join(filepath.Dir(wfqlogPath), PATH_FLAG)
}

func ChipIdFromFlagPath(filePath string) string {
	return strings.Split(filepath.Base(filePath), "_")[0]
}

func HumanReadable(size int64) string {
	div, exp := int64(1024), 0
	for n := size / 1024; n >= 1024; n /= 1024 {
		div *= 1024
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func GetArchiveName(timestamp time.Time) string {
	return timestamp.Format(LAYOUT_ARCHIVE)
}
