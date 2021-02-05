package util

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"
)

const (
	LAYOUT_ARCHIVE = "200601"
)

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

func IsArchiveDir(filePath string) bool {
	r := regexp.MustCompile("^202[0-9]{3}$")
	return r.MatchString(filepath.Base(filePath))
}
