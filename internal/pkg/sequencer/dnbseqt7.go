package sequencer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sharpevo/seqbot/internal/pkg/lane"
	"github.com/sharpevo/seqbot/pkg/util"
)

const (
	PATH_RESULT = "result/OutputFq/upload"
	PATH_FLAG   = "flag"

	ATTR_BARCODE_TYPE = "Barcode Type"
)

type Dnbseqt7 struct{}

func (d *Dnbseqt7) GetBarcode(flagPath string) (string, error) {
	f, err := readFlag(flagPath)
	if err != nil {
		return "", err
	}
	return f.barcodeType(), nil
}

func (d *Dnbseqt7) GetSlide(flagPath string) (string, error) {
	return parseSlideFromFlag(flagPath), nil
}
func (d *Dnbseqt7) GetResultDir(flagPath string) (string, error) {
	return filepath.Join(
		getResultDirFromFlagPath(flagPath),
		parseSlideFromFlag(flagPath),
	), nil
}
func (d *Dnbseqt7) GetArchiveDir(flagPath string) (string, error) {
	return filepath.Join(
		getResultDirFromFlagPath(flagPath),
		util.GetArchiveName(time.Now()),
	), nil
}

func (d *Dnbseqt7) GetWfqTime(flagPath string) (string, error) {
	slide := parseSlideFromFlag(flagPath)
	l := lane.NewLane(slide)
	if err := l.Finish(); err != nil {
		return "", err
	}
	return l.Duration(), nil
}

func (d *Dnbseqt7) GetUploadTime(flagPath string) (string, error) {
	return "-", nil
}

func (d *Dnbseqt7) IsSuccess(flagPath string) (bool, error) {
	return false, nil
}

func parseSlideFromFlag(flagPath string) string {
	return strings.Split(filepath.Base(flagPath), "_")[0]
}

func getResultDirFromFlagPath(flagPath string) string {
	return filepath.Join(
		filepath.Dir(filepath.Dir(filepath.Dir(flagPath))),
		PATH_RESULT)
}

type FlagJson struct {
	ExperimentInfoVec [][]string        `json:"experimentInfoVec"`
	SpeciesBarcodes   map[string]string `json:"speciesBarcodes"`
}

func (f *FlagJson) barcodeType() string {
	for _, e := range f.ExperimentInfoVec {
		if e[0] == ATTR_BARCODE_TYPE {
			return e[1]
		}
	}
	return ""
}

func readFlag(flagPath string) (*FlagJson, error) {
	flagJson := &FlagJson{}
	flagFile, err := os.Open(flagPath)
	if err != nil {
		return flagJson, err
	}
	defer flagFile.Close()
	flagBytes, _ := ioutil.ReadAll(flagFile)
	return flagJson, json.Unmarshal(flagBytes, flagJson)
}
