package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	NAME_BARCODE = "Barcode"

	ATTR_BARCODE_TYPE = "Barcode Type"

	MSG_TPL_BARCODE_SUCC = "**%s**: sequencing completed"
	MSG_BARCODE_FAIL     = "**-**: sequencing completed"
)

type BarcodeAction struct{}

func (b *BarcodeAction) Run(
	eventName string,
	wfqLogPath string,
	chipId string,
) (string, error) {
	f, err := readFlag(eventName)
	if err != nil {
		return MSG_BARCODE_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_BARCODE_SUCC, f.barcodeType()), nil
}

func (b *BarcodeAction) Name() string {
	return NAME_BARCODE
}

type Flag struct {
	ExperimentInfoVec [][]string        `json:"experimentInfoVec"`
	SpeciesBarcodes   map[string]string `json:"speciesBarcodes"`
}

func (f *Flag) barcodeType() string {
	for _, e := range f.ExperimentInfoVec {
		if e[0] == ATTR_BARCODE_TYPE {
			return e[1]
		}
	}
	return ""
}

func readFlag(flagPath string) (*Flag, error) {
	flag := &Flag{}
	flagFile, err := os.Open(flagPath)
	if err != nil {
		return flag, err
	}
	defer flagFile.Close()
	flagBytes, _ := ioutil.ReadAll(flagFile)
	return flag, json.Unmarshal(flagBytes, flag)
}
