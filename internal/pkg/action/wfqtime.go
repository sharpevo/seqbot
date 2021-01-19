package action

import (
	"fmt"

	"github.com/sharpevo/seqbot/internal/pkg/lane"
)

const (
	NAME_WFQTIME = "WfqTime"

	MSG_TPL_WFQTIME_SUCC = "- WFQ Time: %s"
	MSG_TPL_WFQTIME_FAIL = "- WFQ Time: ?"
)

type WfqTimeAction struct{}

func (w *WfqTimeAction) Run(
	eventName string,
	wfqLogPath string,
	chipId string,
) (string, error) {
	l := lane.NewLane(chipId)
	if err := l.Finish(); err != nil {
		return MSG_TPL_WFQTIME_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_WFQTIME_SUCC, l.Duration()), nil
}

func (w *WfqTimeAction) Name() string {
	return NAME_WFQTIME
}
