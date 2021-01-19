package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sharpevo/seqbot/pkg/util"
)

const (
	NAME_SUMMARY = "Summary"

	MSG_TPL_SUMMARY_SUCC = "- Count: %d\n- Size: %s"
	MSG_SUMMARY_FAIL     = "- Count: -\n- Size: -"
)

type SummaryAction struct{}

func (s *SummaryAction) Run(
	eventName string,
	wfqLogPath string,
	chipId string,
) (string, error) {
	resultChipPath := util.ResultChipPathFromWFQLogPath(wfqLogPath, chipId)
	var size int64
	count := 0
	err := filepath.Walk(resultChipPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".fq.gz") {
			count++
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return MSG_SUMMARY_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_SUMMARY_SUCC, count, util.HumanReadable(size)), nil
}

func (s *SummaryAction) Name() string {
	return NAME_SUMMARY
}
