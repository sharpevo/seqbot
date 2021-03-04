package sequencer

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sharpevo/seqbot/pkg/util"
)

const (
	FILE_RUN_PARAMETERS = "RunParameters.xml"
	FILE_RUN_COMPLETION = "RunCompletionStatus.xml"
)

type Nextseq500 struct{}

func (n *Nextseq500) GetBarcode(successPath string) (string, error) {
	parameterPath := filepath.Join(filepath.Dir(successPath), FILE_RUN_PARAMETERS)
	runParameter, err := readParameter(parameterPath)
	if err != nil {
		return "", err
	}
	return runParameter.LibraryId, nil
}

func (n *Nextseq500) GetSlide(successPath string) (string, error) {
	return filepath.Base(filepath.Dir(successPath)), nil
}

func (n *Nextseq500) GetArchiveDir(successPath string) (string, error) {
	return filepath.Join(
		filepath.Dir(filepath.Dir(successPath)),
		util.GetArchiveName(time.Now()),
	), nil
	return "", fmt.Errorf("archive not supported for NextSeq 500")
}

func (n *Nextseq500) GetResultDir(successPath string) (string, error) {
	return filepath.Dir(successPath), nil
}

func (n *Nextseq500) GetWfqTime(successPath string) (string, error) {
	return "", fmt.Errorf("wfqtime not supported for NextSeq 500")
}

func (n *Nextseq500) GetUploadTime(successPath string) (string, error) {
	return "", fmt.Errorf("uploadtime not supported for NextSeq 500")
}

func (n *Nextseq500) IsSuccess(filePath string) (bool, error) {
	fileName := filepath.Base(filePath)
	return fileName == FILE_RUN_COMPLETION, nil
}

type RunParameter struct {
	LibraryId string `xml:"LibraryID"`
	RunId     string `xml:"RunID"`
}

func readParameter(flagPath string) (*RunParameter, error) {
	runParameter := &RunParameter{}
	file, err := os.Open(flagPath)
	if err != nil {
		return runParameter, err
	}
	defer file.Close()
	runParameterBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return runParameter, err
	}
	return runParameter, xml.Unmarshal(runParameterBytes, runParameter)
}
