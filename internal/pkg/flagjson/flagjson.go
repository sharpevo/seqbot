package flagjson

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Flag struct {
	ExperimentInfoVec [][]string        `json:"experimentInfoVec"`
	SpeciesBarcodes   map[string]string `json:"speciesBarcodes"`
}

func (f *Flag) BarcodeType() string {
	for _, e := range f.ExperimentInfoVec {
		if e[0] == "Barcode Type" {
			return e[1]
		}
	}
	return ""
}

func ReadFlag(flagPath string) (*Flag, error) {
	flag := &Flag{}
	flagFile, err := os.Open(flagPath)
	if err != nil {
		return flag, err
	}
	defer flagFile.Close()
	flagBytes, _ := ioutil.ReadAll(flagFile)
	return flag, json.Unmarshal(flagBytes, flag)
}
