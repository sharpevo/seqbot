package util

import (
	"fmt"
	"os"
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

func FastqCountAndSize(fastqPath string) (int, string, error) {
	var size int64
	count := 0
	err := filepath.Walk(fastqPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".fq.gz") {
			count++
			size += info.Size()
		}
		return err
	})
	return count, humanReadable(size), err
}

func humanReadable(size int64) string {
	div, exp := int64(1024), 0
	for n := size / 1024; n >= 1024; n /= 1024 {
		div *= 1024
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func archivedDir(timestamp time.Time) string {
	return timestamp.Format(LAYOUT_ARCHIVE)
}

func CreateArchivedDir(rootDir string, timestamp time.Time) (string, error) {
	archivedPath := filepath.Join(rootDir, archivedDir(timestamp))
	if err := os.MkdirAll(archivedPath, 0755); err != nil {
		return archivedPath, err
	}
	return archivedPath, nil
}
