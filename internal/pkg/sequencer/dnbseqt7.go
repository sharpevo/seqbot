package sequencer

import (
	"encoding/json"
	"fmt"
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

	ATTR_BARCODE_TYPE  = "Barcode Type"
	ATTR_SEQUENCE_TYPE = "Sequence Type"
	ATTR_DUAL_BARCODE  = "Dual Barcode"

	VALUE_BARCODE_SINGLE = "single"
	VALUE_BARCODE_DUAL   = "dual"
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

func (d *Dnbseqt7) GetExtraExperimentInfo(flagPath string) (string, error) {
	f, err := readFlag(flagPath)
	if err != nil {
		return "", err
	}
	sequenceType := f.sequenceType()
	dualBarcode := VALUE_BARCODE_SINGLE
	if f.dualBarcode() != "0" {
		dualBarcode = VALUE_BARCODE_DUAL
	}
	return fmt.Sprintf("%s, %s", sequenceType, dualBarcode), nil
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
	info, err := os.Stat(flagPath)
	if err != nil {
		return "", nil
	}
	latest := info.ModTime()
	earliest := latest
	resultRootPath := getResultDirFromFlagPath(flagPath)
	resultPath := filepath.Join(resultRootPath, parseSlideFromFlag(flagPath))
	err = filepath.Walk(
		resultPath,
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
	return f.get(ATTR_BARCODE_TYPE)
}

func (f *FlagJson) sequenceType() string {
	return f.get(ATTR_SEQUENCE_TYPE)
}

func (f *FlagJson) dualBarcode() string {
	return f.get(ATTR_DUAL_BARCODE)
}

func (f *FlagJson) get(attrName string) string {
	for _, e := range f.ExperimentInfoVec {
		if e[0] == attrName {
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
