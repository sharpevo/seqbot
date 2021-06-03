package sequencer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sharpevo/seqbot/pkg/util"

	"github.com/sirupsen/logrus"
)

const (
	BIOINFO_NAME = "BioInfo.csv"
	PREFIX_DNBID = "DNB ID,"

	TMPL_SUCCESS = "^%s.*_Success.txt$"
)

var lanes = []string{"L01", "L02", "L03", "L04"}

type Mgiseq2000 struct{}

func (m *Mgiseq2000) GetBarcode(successPath string) (string, error) {
	dir := filepath.Dir(successPath)
	idMap := map[string]struct{}{}
	ids := []string{}
	for _, l := range lanes {
		bioinfoPath := filepath.Join(dir, l, BIOINFO_NAME)
		bioinfoFile, err := os.Open(bioinfoPath)
		if err != nil {
			logrus.Warnf("cannot parse dnb id from %s: %v", bioinfoPath, err)
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
	return strings.Join(ids, ","), nil
}

func (m *Mgiseq2000) GetSlide(successPath string) (string, error) {
	return filepath.Base(filepath.Dir(successPath)), nil
}

func (m *Mgiseq2000) GetExtraExperimentInfo(successPath string) (string, error) {
	return "", fmt.Errorf("extra info for mgi2000 has not been implemented")
}

func (m *Mgiseq2000) GetArchiveDir(successPath string) (string, error) {
	return filepath.Join(
		filepath.Dir(filepath.Dir(successPath)),
		util.GetArchiveName(time.Now()),
	), nil
}

func (m *Mgiseq2000) GetResultDir(successPath string) (string, error) {
	return filepath.Dir(successPath), nil
}

func (m *Mgiseq2000) GetWfqTime(successPath string) (string, error) {
	return "", nil
}

func (m *Mgiseq2000) GetUploadTime(successPath string) (string, error) {
	info, err := os.Stat(successPath)
	if err != nil {
		return "", err
	}
	latest := info.ModTime()
	earliest := latest
	err = filepath.Walk(
		filepath.Dir(successPath),
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.ModTime().Before(earliest) {
				earliest = info.ModTime()
			}
			return nil
		},
	)
	return fmt.Sprintf("%v", latest.Sub(earliest).Round(time.Second)), err
}

func (m *Mgiseq2000) IsSuccess(filePath string) (bool, error) {
	lastDir := filepath.Base(filepath.Dir(filePath))
	r := regexp.MustCompile(fmt.Sprintf(TMPL_SUCCESS, lastDir))
	return r.MatchString(filepath.Base(filePath)), nil
}
