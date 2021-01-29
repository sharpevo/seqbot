package action

import (
	"fmt"
)

const (
	NAME_WFQTIME = "WfqTime"

	MSG_TPL_WFQTIME_SUCC = "- WFQ Time: %s"
	MSG_WFQTIME_FAIL     = "- WFQ Time: -"
)

type WfqTimeAction struct{}

func (w *WfqTimeAction) Run(
	eventName string,
	command CommandInterface,
) (string, error) {
	wfqtime, err := command.Sequencer().GetWfqTime(eventName)
	if err != nil {
		return MSG_WFQTIME_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_WFQTIME_SUCC, wfqtime), nil
}

func (w *WfqTimeAction) Name() string {
	return NAME_WFQTIME
}
