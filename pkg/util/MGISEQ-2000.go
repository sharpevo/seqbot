package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	BIOINFO_NAME = "BioInfo.csv"
	PREFIX_DNBID = "DNB ID,"

	TMPL_SUCCESS = "^%s.*_Success.txt$"
)

var lanes = []string{"L01", "L02", "L03", "L04"}

func IsSuccess(filePath string) bool {
	lastDir := filepath.Base(filepath.Dir(filePath))
	r := regexp.MustCompile(fmt.Sprintf(TMPL_SUCCESS, lastDir))
	return r.MatchString(filepath.Base(filePath))
}

func ParseMgiInfo(successPath string) (string, string) {
	dir := filepath.Dir(successPath)
	idMap := map[string]struct{}{}
	ids := []string{}
	for _, l := range lanes {
		bioinfoPath := filepath.Join(dir, l, BIOINFO_NAME)
		bioinfoFile, err := os.Open(bioinfoPath)
		if err != nil {
			logrus.Warnf("cannot parse dnb id from %s", bioinfoPath)
			continue
		}
		defer bioinfoFile.Close()
		scanner := bufio.NewScanner(bioinfoFile)
		for scanner.Scan() {
			splits := strings.Split(scanner.Text(), PREFIX_DNBID)
			if len(splits) == 2 {
				if _, ok := idMap[splits[1]]; !ok {
					ids = append(ids, splits[1])
					idMap[splits[1]] = struct{}{}
				}
				break
			}
		}
	}
	return filepath.Base(dir), strings.Join(ids, ",")
}
